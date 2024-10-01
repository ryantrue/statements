package config

import (
	"github.com/spf13/viper"
)

// Config структура для хранения конфигурации
type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		URL                string `mapstructure:"url"`
		MaxConnections     int    `mapstructure:"max_connections"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections"`
		MigrationsDir      string `mapstructure:"migrations_dir"`
		DatabaseName       string `mapstructure:"database_name"`
	} `mapstructure:"database"`

	FileUpload struct {
		UploadDir string `mapstructure:"upload_dir"`
		StaticDir string `mapstructure:"static_dir"`
	} `mapstructure:"file_upload"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	Python struct {
		Interpreter string `mapstructure:"interpreter"`
		ScriptPath  string `mapstructure:"script_path"`
	} `mapstructure:"python"`
}

// LoadConfig загружает конфигурацию с помощью Viper
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv() // Чтение переменных окружения

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
