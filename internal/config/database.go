package config

import "github.com/spf13/viper"

// DatabaseConfig конфигурация для базы данных
type DatabaseConfig struct {
	URL                string `mapstructure:"url"`
	MaxConnections     int    `mapstructure:"max_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	MigrationsDir      string `mapstructure:"migrations_dir"`
	DatabaseName       string `mapstructure:"database_name"`
}

// LoadDatabaseConfig загружает конфигурацию базы данных
func LoadDatabaseConfig(v *viper.Viper) (DatabaseConfig, error) {
	var config DatabaseConfig
	if err := v.UnmarshalKey("database", &config); err != nil {
		return config, err
	}
	return config, nil
}
