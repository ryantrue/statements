package config

import "github.com/spf13/viper"

// LoggingConfig конфигурация для логирования
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// LoadLoggingConfig загружает конфигурацию логирования
func LoadLoggingConfig(v *viper.Viper) (LoggingConfig, error) {
	var config LoggingConfig
	if err := v.UnmarshalKey("logging", &config); err != nil {
		return config, err
	}
	return config, nil
}
