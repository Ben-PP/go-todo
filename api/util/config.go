package util

import "github.com/spf13/viper"

type Config struct {
    DbDriver         string `mapstructure:"DB_DRIVER"`
    DbUrl         string `mapstructure:"DB_URL"`
    PostgresUser     string `mapstructure:"POSTGRES_USER"`
    PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
    PostgresDb       string `mapstructure:"POSTGRES_DB"`
    ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
	// TODO Choose programmatically the env to load
    viper.SetConfigName("dev")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}