package minioutils

import (
	"strings"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
)

type MinioConfig struct {
	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

func ProvideMinioConfig(cl *utils.ConfigLocation) (*MinioConfig, error) {
	cfg := &MinioConfig{}
	viper.AddConfigPath(cl.ConfigPath)
	viper.SetConfigName(cl.ConfigName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	isNotFoundError := func(m string) bool {
		return strings.Contains(strings.ToLower(m), "not found")
	}
	err := viper.ReadInConfig()
	if err != nil && !isNotFoundError(err.Error()) {
		return cfg, err
	}
	if err != nil && isNotFoundError(err.Error()) {
		cfg.MinioEndpoint = viper.GetString("MINIO_ENDPOINT")
		cfg.MinioAccessKey = viper.GetString("MINIO_ACCESS_KEY")
		cfg.MinioSecretKey = viper.GetString("MINIO_SECRET_KEY")
		cfg.MinioUseSSL = viper.GetBool("MINIO_USE_SSL")

		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}
