package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type ApiConfig struct {
	errors []error
	Env    string
	Port   int
}

func LoadApiConfig() ApiConfig {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()
	c := ApiConfig{}
	c.Env = c.getMandatoryString("ENV")
	c.Port = c.getIntWithDefault("PORT", 8080)

	if len(c.errors) != 0 {
		errorReport := "errors in config :\n"
		for _, err := range c.errors {
			errorReport += fmt.Sprintf("- %s\n", err)
		}
		panic(fmt.Errorf(errorReport))
	}
	return c
}

func (c *ApiConfig) getIntWithDefault(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}

func (c *ApiConfig) getMandatoryString(key string) (value string) {
	if value = viper.GetString(key); value == "" {
		c.errors = append(c.errors, fmt.Errorf("cannot find configuration for key %s", key))
	}
	return value
}
