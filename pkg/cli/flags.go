package cli

type SymbioticSchema struct {
	SymbioticApiUrl string `help:"Symbiotic API URL." default:"https://app.symbiotic.fi/api/v2/vaults" name:"symbiotic.api"`
}

type EthereumSchema struct {
	EthereumRPC string `help:"Ethereum RPC URL."  default:"http://127.0.0.1:8545" name:"ethereum.rpc"`
	EthereumWS  string `help:"Ethereum WebSocket URL."  default:"http://127.0.0.1:8546" name:"ethereum.ws"`
}

type PrometheusSchema struct {
	PrometheusHost string `help:"Prometheus host address." default:"0.0.0.0" name:"prometheus.host"`
	PrometheusPort int    `help:"Prometheus port number." default:"9090" name:"prometheus.port"`
}

type ConfigSchema struct {
	ConfigFile string `help:"[OPTIONAL] Path to the configuration file (toml) overrides all the flags " default:"" name:"config.file"`
}
