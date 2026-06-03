package ai

import (
	"context"
	"fmt"
	"log"

	"logic-grpc-service/config"
	"logic-grpc-service/internal/ai/tools"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

type ChatAgent struct {
	agent     *react.Agent
	chatModel *openai.ChatModel
}

var chatAgent *ChatAgent

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	HrId    int64  `json:"hr_id,omitempty"`
}

func Init(db *gorm.DB) error {
	cfg := config.GetAIConfig()

	if cfg.APIKey == "" || cfg.Endpoint == "" {
		return fmt.Errorf("AI配置不完整，请检查api_key和endpoint")
	}

	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.Endpoint,
	})
	if err != nil {
		return fmt.Errorf("初始化ChatModel失败: %w", err)
	}

	allTools := []tool.BaseTool{
		tools.NewQueryApplicationStatsTool(db),
		tools.NewQueryPositionStatsTool(db),
		tools.NewQueryCandidatesTool(db),
		tools.NewQueryPositionHotnessTool(db),
	}

	toolsConfig := compose.ToolsNodeConfig{
		Tools: allTools,
	}

	agent, err := react.NewAgent(context.Background(), &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      toolsConfig,
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			systemMsg := schema.SystemMessage(
				`你是智能招聘系统的AI助手。你可以帮助HR用户：
1. 查询招聘统计数据（投递总人数、各状态分布）
2. 查询单岗位投递统计（投递量、学历分布、技能分析）
3. 筛选符合条件的候选人（按学历、技能、经验筛选）
4. 查看岗位热度排行
5. 提供准确的数据，不要编撰数据，避免使用错误的统计数据。
6. 若无足够信息，请回答"暂无相关信息"并说明原因。

 请根据用户的自然语言请求，调用合适的工具完成任务。回答时请使用中文，语气友好专业。`)
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, systemMsg)
			res = append(res, input...)
			return res
		},
		MaxStep: 10,
	})
	if err != nil {
		return fmt.Errorf("创建Agent失败: %w", err)
	}

	chatAgent = &ChatAgent{agent: agent, chatModel: chatModel}
	return nil
}

func GetChatAgent() *ChatAgent {
	return chatAgent
}

func (ca *ChatAgent) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	if ca == nil || ca.agent == nil {
		return "", fmt.Errorf("AI客户端未初始化")
	}

	schemaMessages := make([]*schema.Message, 0, len(messages)+1)

	for _, msg := range messages {
		if msg.HrId > 0 {
			schemaMessages = append(schemaMessages, schema.SystemMessage(
				fmt.Sprintf("当前HR用户的ID为: %d。调用任何工具时，请使用此ID作为hr_id参数。", msg.HrId),
			))
			break
		}
	}

	for _, msg := range messages {
		if msg.Content == "" {
			continue
		}
		switch msg.Role {
		case "user":
			schemaMessages = append(schemaMessages, schema.UserMessage(msg.Content))
		case "assistant":
			if msg.Content != "" {
				schemaMessages = append(schemaMessages, schema.AssistantMessage(msg.Content, []schema.ToolCall{}))
			}
		case "system":
			schemaMessages = append(schemaMessages, schema.SystemMessage(msg.Content))
		}
	}

	resp, err := ca.agent.Generate(ctx, schemaMessages)
	if err != nil {
		log.Printf("AI生成回复失败: %v", err)
		return "", fmt.Errorf("服务暂不可用，请稍后再试")
	}

	return resp.Content, nil
}

func (ca *ChatAgent) ChatStream(ctx context.Context, messages []ChatMessage) (<-chan string, error) {
	if ca == nil || ca.agent == nil || ca.chatModel == nil {
		return nil, fmt.Errorf("AI客户端未初始化")
	}

	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("AI流式响应panic: %v", r)
			}
		}()

		schemaMessages := make([]*schema.Message, 0, len(messages)+1)
		for _, msg := range messages {
			if msg.HrId > 0 {
				schemaMessages = append(schemaMessages, schema.SystemMessage(
					fmt.Sprintf("当前HR用户的ID为: %d。调用任何工具时，请使用此ID作为hr_id参数。", msg.HrId),
				))
				break
			}
		}
		for _, msg := range messages {
			if msg.Content == "" {
				continue
			}
			switch msg.Role {
			case "user":
				schemaMessages = append(schemaMessages, schema.UserMessage(msg.Content))
			case "assistant":
				if msg.Content != "" {
					schemaMessages = append(schemaMessages, schema.AssistantMessage(msg.Content, []schema.ToolCall{}))
				}
			case "system":
				schemaMessages = append(schemaMessages, schema.SystemMessage(msg.Content))
			}
		}

		agentResp, err := ca.agent.Generate(ctx, schemaMessages)
		if err != nil {
			log.Printf("AI生成回复失败: %v", err)
			ch <- "服务暂不可用，请稍后再试"
			return
		}

		finalMessages := append(schemaMessages, schema.AssistantMessage(agentResp.Content, []schema.ToolCall{}))

		stream, err := ca.chatModel.Stream(ctx, finalMessages)
		if err != nil {
			log.Printf("流式生成失败: %v", err)
			ch <- "服务暂不可用，请稍后再试"
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				chunk, err := stream.Recv()
				if err != nil {
					// 流式结束或出错
					return
				}
				ch <- chunk.Content
			}
		}
	}()

	return ch, nil
}

func (ca *ChatAgent) ChatWithPrompt(ctx context.Context, prompt string) (string, error) {
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}
	return ca.Chat(ctx, messages)
}

func Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	if chatAgent == nil {
		return "", fmt.Errorf("AI客户端未初始化")
	}
	return chatAgent.Chat(ctx, messages)
}

func ChatStream(ctx context.Context, messages []ChatMessage) (<-chan string, error) {
	if chatAgent == nil {
		return nil, fmt.Errorf("AI客户端未初始化")
	}
	return chatAgent.ChatStream(ctx, messages)
}

func ChatWithPrompt(ctx context.Context, prompt string) (string, error) {
	if chatAgent == nil {
		return "", fmt.Errorf("AI客户端未初始化")
	}
	return chatAgent.ChatWithPrompt(ctx, prompt)
}
