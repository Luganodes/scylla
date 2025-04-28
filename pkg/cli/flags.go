package cli

type SymbioticSchema struct {
	SymbioticApiUrl string `help:"Symbiotic API URL." default:"https://app.symbiotic.fi/api/v2/vaults" name:"symbiotic.api"`
}

type EthereumSchema struct {
	EthereumRPC string `help:"Ethereum RPC URL."  default:"127.0.0.1" name:"ethereum.rpc"`
	EthereumWS  string `help:"Ethereum WebSocket URL."  default:"127.0.0.1" name:"ethereum.ws"`
}

type PrometheusSchema struct {
	PrometheusHost string `help:"Prometheus host address." default:"0.0.0.0" name:"prometheus.host"`
	PrometheusPort int    `help:"Prometheus port number." default:"9090" name:"prometheus.port"`
}
