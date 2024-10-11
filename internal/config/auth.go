package config

import "github.com/spf13/viper"

// AuthConfig конфигурация для авторизации (JWT)
type AuthConfig struct {
	JWTSecret string `mapstructure:"jwtSecret"`
}

// LoadAuthConfig загружает конфигурацию авторизации (JWT)
func LoadAuthConfig(v *viper.Viper) (AuthConfig, error) {
	var config AuthConfig
	if err := v.UnmarshalKey("auth", &config); err != nil {
		return config, err
	}
	return config, nil
}
