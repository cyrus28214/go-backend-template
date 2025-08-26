package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
	Mode    string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN        string           `mapstructure:"dsn"`
	GormLogger GormLoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level     string `mapstructure:"level"`
	Dir       string `mapstructure:"dir"`
	AddSource bool   `mapstructure:"add_source"`
}

type GormLoggerConfig struct {
	Level         string `mapstructure:"level"`
	SlowThreshold int    `mapstructure:"slow_threshold"`
}

type JwtConfig struct {
	JwtSecretHex string `mapstructure:"jwt_secret_hex"`
}

func LoadConfig(path string) *Config {
	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &cfg
}
