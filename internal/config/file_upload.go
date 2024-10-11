package config

import "github.com/spf13/viper"

// FileUploadConfig конфигурация для загрузки файлов
type FileUploadConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
	StaticDir string `mapstructure:"static_dir"`
}

// LoadFileUploadConfig загружает конфигурацию загрузки файлов
func LoadFileUploadConfig(v *viper.Viper) (FileUploadConfig, error) {
	var config FileUploadConfig
	if err := v.UnmarshalKey("file_upload", &config); err != nil {
		return config, err
	}
	return config, nil
}
