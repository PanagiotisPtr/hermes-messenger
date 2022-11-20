package minioutils

import (
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
)

type MinioConfig struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

func ProvideMinioConfig(cl *utils.ConfigLocation) (*MinioConfig, error) {
	cfg := &MinioConfig{}
	viper.AddConfigPath(cl.ConfigPath)
	viper.SetConfigName(cl.ConfigName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}
