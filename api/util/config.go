package util

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
    DbUrl                   string  `mapstructure:"DB_URL"`
    AuthTokenLifeSpan       int     `mapstructure:"AUTH_TOKEN_LIFE_SPAN"`
    RefreshTokenLifeSpan    int     `mapstructure:"REFRESH_TOKEN_LIFE_SPAN"`
    JwtSecret               string  `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)

    if os.Getenv("GO_ENV") == "dev" {
        viper.SetConfigName("dev")
    } else {
        viper.SetConfigName("prod")
    }
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}