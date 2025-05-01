// config/config.go
package config

import (
	"log"
	"path/filepath"

	"github.com/luganodes/slashing-observer/common/utils"
	"github.com/spf13/viper"
)

var (
	EXTERNAL_APP_NAME string
	EXTERNAL_VERSION  string
)

var (
	API_URL string
	RPC_URL string
	WS_URL  string

	PROMETHEUS_HOST string
	PROMETHEUS_PORT int
)

func LoadConfig(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln("[ERROR] Failed to get absolute path:", err)
	}
	utils.LoadConfig(absPath)
	API_URL = viper.GetString("symbiotic.api_url")
	RPC_URL = viper.GetString("ethereum.rpc_url")
	WS_URL = viper.GetString("ethereum.ws_url")
	PROMETHEUS_HOST = viper.GetString("prometheus.host")
	PROMETHEUS_PORT = viper.GetInt("prometheus.port")
}

func SetSymbiotiURL(url string) {
	API_URL = url
}

func SetEthereumURL(rpcurl string, wsurl string) {
	RPC_URL = rpcurl
	WS_URL = wsurl
}

func SetprometheusData(host string, port int) {
	PROMETHEUS_HOST = host
	PROMETHEUS_PORT = port
}
