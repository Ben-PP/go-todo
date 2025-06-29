package util

import (
	"errors"
	gterrors "go-todo/gt_errors"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
    DbUrl                   string  `mapstructure:"DB_URL"`
    AccessTokenLifeSpan     int     `mapstructure:"ACCESS_TOKEN_LIFE_SPAN"`
    RefreshTokenLifeSpan    int     `mapstructure:"REFRESH_TOKEN_LIFE_SPAN"`
    JwtAccessSecret         string  `mapstructure:"JWT_ACCESS_SECRET"`
    JwtRefreshSecret        string  `mapstructure:"JWT_REFRESH_SECRET"`
}

var globalConfig *Config

func LoadConfig(path string) (config *Config, err error) {
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

func GetConfig() (config *Config, err error) {
    if globalConfig == nil {
        globalConfig, err = LoadConfig(".")
        if err != nil {
            err = errors.Join(gterrors.ErrConfigLoadFailed, err)
            return
        }
    }
    config = globalConfig
    return
}