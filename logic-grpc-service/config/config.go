package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OSS      OSSConfig
	AI       AIConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

type AIConfig struct {
	Provider    string
	Model       string
	APIKey      string
	Endpoint    string
	Temperature float64
}

type JWTConfig struct {
	Secret string
	Expire int
}

var cfg *Config

func Init() error {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg = &Config{
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetString("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			Name:     v.GetString("DB_NAME"),
		},
		OSS: OSSConfig{
			Endpoint:        v.GetString("OSS_ENDPOINT"),
			AccessKeyID:     v.GetString("OSS_ACCESS_KEY_ID"),
			AccessKeySecret: v.GetString("OSS_ACCESS_KEY_SECRET"),
			BucketName:      v.GetString("OSS_BUCKET_NAME"),
		},
		AI: AIConfig{
			Provider:    v.GetString("AI_PROVIDER"),
			Model:       v.GetString("AI_MODEL"),
			APIKey:      v.GetString("AI_API_KEY"),
			Endpoint:    v.GetString("AI_ENDPOINT"),
			Temperature: v.GetFloat64("AI_TEMPERATURE"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
			Expire: v.GetInt("JWT_EXPIRE"),
		},
	}

	return nil
}

func Get() *Config {
	return cfg
}

func GetServerConfig() ServerConfig {
	return cfg.Server
}

func GetDatabaseConfig() DatabaseConfig {
	return cfg.Database
}

func GetOSSConfig() OSSConfig {
	return cfg.OSS
}

func GetAIConfig() AIConfig {
	return cfg.AI
}

func GetJWTConfig() JWTConfig {
	return cfg.JWT
}
