// config/config.go
package config

import (
	"github.com/luganodes/slashing-observer/common/utils"
	"github.com/spf13/viper"
)

var (
	API_URL string
	RPC_URL string
	WS_URL  string
)

func LoadConfig(path string) {
	utils.LoadConfig(path)
	API_URL = viper.GetString("symbiotic.api_url")
	RPC_URL = viper.GetString("ethereum.rpc_url")
	WS_URL = viper.GetString("ethereum.ws_url")
}
