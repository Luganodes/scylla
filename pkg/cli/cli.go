package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/luganodes/slashing-observer/config"
)

var CLIStruct struct {
	Version kong.VersionFlag `help:"Show app version." name:"version"`

	EthereumSchema
	PrometheusSchema
	SymbioticSchema
	ConfigSchema

	Start StartCmd `cmd:"" help:"Helps to start the Slashing observer"`
}

func RunCli() {
	ctx := kong.Parse(&CLIStruct,
		kong.Name(config.EXTERNAL_APP_NAME),
		kong.Description("Symbiotic Slashing Observer"),
		kong.Vars{"version": config.EXTERNAL_VERSION})

	CLIStruct.ConfigFile = strings.TrimSpace(CLIStruct.ConfigFile)
	configFile := CLIStruct.ConfigFile != "" && strings.HasSuffix(CLIStruct.ConfigFile, ".toml")

	if configFile {
		config.LoadConfig(CLIStruct.ConfigFile)
	} else {
		config.SetEthereumURL(
			CLIStruct.EthereumRPC,
			CLIStruct.EthereumWS)
		config.SetSymbiotiURL(
			CLIStruct.SymbioticApiUrl,
		)
		config.SetprometheusData(
			CLIStruct.PrometheusHost,
			CLIStruct.PrometheusPort,
		)
	}
	fmt.Println(config.PROMETHEUS_HOST, config.PROMETHEUS_PORT)
	fmt.Println(config.API_URL)
	fmt.Println(config.WS_URL, config.RPC_URL)
	if err := ctx.Run(&CLIStruct); err != nil {
		log.Fatalln(err)
	}
}
