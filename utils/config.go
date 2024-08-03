package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Hold all env var with viper
type Configuration struct {
	DbDriver string `mapstructure:"DB_DRIVER"`
	DbSource string `mapstructure:"DB_SOURCE"`
	HTTPSServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration 	`mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json,xml works

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	
	err = viper.Unmarshal(&config)
	return
}