package config

import "github.com/spf13/viper"

// ServerConfig конфигурация для сервера
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LoadServerConfig загружает конфигурацию сервера
func LoadServerConfig(v *viper.Viper) (ServerConfig, error) {
	var config ServerConfig
	if err := v.UnmarshalKey("server", &config); err != nil {
		return config, err
	}
	return config, nil
}
