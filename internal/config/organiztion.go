package config

import "github.com/spf13/viper"

// OrganizationConfig конфигурация для организации
type OrganizationConfig struct {
	DefaultInn        string `mapstructure:"default_inn"`
	DefaultName       string `mapstructure:"default_name"`
	DefaultInnCredit  string `mapstructure:"default_inn_credit"`
	DefaultNameCredit string `mapstructure:"default_name_credit"`
}

// LoadOrganizationConfig загружает конфигурацию организации
func LoadOrganizationConfig(v *viper.Viper) (OrganizationConfig, error) {
	var config OrganizationConfig
	if err := v.UnmarshalKey("organization", &config); err != nil {
		return config, err
	}
	return config, nil
}
