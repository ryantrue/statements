package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Config объединяет все конфигурационные секции
type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	FileUpload   FileUploadConfig   `mapstructure:"file_upload"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Python       PythonConfig       `mapstructure:"python"`
	Auth         AuthConfig         `mapstructure:"auth"`
	Organization OrganizationConfig `mapstructure:"organization"`
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig(configPath string) (*Config, error) {
	v := initViper(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Можно добавить валидацию параметров здесь
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// initViper инициализирует Viper с настройками для работы с окружением и конфигурационным файлом
func initViper(configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	return v
}

// validateConfig выполняет проверку критических параметров конфигурации
func validateConfig(config *Config) error {
	// Пример проверки значений
	if config.Server.Host == "" {
		return fmt.Errorf("server host is not set")
	}
	if config.Server.Port == 0 {
		return fmt.Errorf("server port is not set")
	}
	// Можно добавить другие проверки для важных параметров
	return nil
}
