package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/luganodes/slashing-observer/config"
)

var CLIStruct struct {
	Version bool `help:"Show app version." name:"version"`

	EthereumSchema
	PrometheusSchema
	SymbioticSchema

	Start StartCmd `cmd:"" help:"Helps to start the Slashing observer"`
}

func RunCli() {
	ctx := kong.Parse(&CLIStruct,
		kong.Name(config.EXTERNAL_APP_NAME),
		kong.Description("Symbiotic Slashing Observer"))

	if CLIStruct.Version {
		fmt.Printf("App version: %s\n", config.EXTERNAL_VERSION)
		os.Exit(0)
	}

	if err := ctx.Run(&CLIStruct); err != nil {
		log.Fatalln(err)
	}
}
