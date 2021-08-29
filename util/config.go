package util

import (
	"github.com/spf13/viper"
)

// Config store all configuration of application
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"Server_Address"`
}

// In order to get the value of the variables and store them in this struct,
// we need to use the unmarshaling feature of Viper.
// viper use mapstructure package
//LoadConfig() which takes a path as input, and returns a config object or an error.
//this function will read configurations from a config file inside the path if it exists,
//or override their values with environment variables if they’re provided.

func LoadConfig(path string) (config Config, err error) {
	// find the config file
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	// start to read config file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	// convert the config to struct
	err = viper.Unmarshal(&config)
	return
}
