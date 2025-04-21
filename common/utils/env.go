package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}
}
