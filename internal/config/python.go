package config

import "github.com/spf13/viper"

// PythonConfig конфигурация для Python
type PythonConfig struct {
	Interpreter string `mapstructure:"interpreter"`
	ScriptPath  string `mapstructure:"script_path"`
}

// LoadPythonConfig загружает конфигурацию для Python
func LoadPythonConfig(v *viper.Viper) (PythonConfig, error) {
	var config PythonConfig
	if err := v.UnmarshalKey("python", &config); err != nil {
		return config, err
	}
	return config, nil
}
