package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	GRPC   GRPCConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
	Expire string
}

type GRPCConfig struct {
	Address string
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
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
			Expire: v.GetString("JWT_EXPIRE"),
		},
		GRPC: GRPCConfig{
			Address: v.GetString("GRPC_ADDR"),
		},
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Get() *Config {
	return cfg
}

func GetServerConfig() ServerConfig {
	return cfg.Server
}

func GetJWTConfig() JWTConfig {
	return cfg.JWT
}

func GetGRPCAddress() string {
	return cfg.GRPC.Address
}
