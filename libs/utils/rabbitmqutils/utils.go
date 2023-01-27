package rabbitmqutils

import (
	"context"
	"strings"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AmqpConfig struct {
	AMQPURI string `mapstructure:"AMQP_URI"`
}

// ProvideAmqpChannel provides an amqp channel
func ProvideAmqpChannel(
	lc fx.Lifecycle,
	cfg *AmqpConfig,
) (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			channelCloseErr := ch.Close()
			connCloseErr := conn.Close()
			if channelCloseErr != nil {
				return channelCloseErr
			}

			return connCloseErr
		},
	})

	return ch, nil
}

func ProvideAmqpConfig(cl *utils.ConfigLocation) (*AmqpConfig, error) {
	cfg := &AmqpConfig{}
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
		cfg.AMQPURI = viper.GetString("AMQP_URI")

		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}
