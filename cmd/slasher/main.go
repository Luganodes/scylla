package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/luganodes/slashing-observer/config"
	"github.com/luganodes/slashing-observer/pkg/observer"
	"github.com/luganodes/slashing-observer/pkg/vault"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "Path to the config file")
	flag.Parse()

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalln("[ERROR] Failed to get absolute path:", err)
	}

	config.LoadConfig(absPath)
}

func main() {
	vaults, err := vault.GetVaultInfoList()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, eachvault := range vaults {
		if eachvault.Slashable {
			slasher, err := vault.GetSlasherAddress(config.RPC_URL, eachvault.Address)
			if err != nil {
				log.Println("Error getting slasher address:", err)
				continue
			}

			slasherType, err := vault.GetSlasherType(config.RPC_URL, slasher)
			if err != nil {
				log.Println("Error getting slasher TYPE:", err)
				continue
			}

			log.Printf("✅ Vault detials: \n	 Vault: %s \n	 Slasher: %s \n	 Type: %d \n	 Name: %s", eachvault.Address, slasher, slasherType, eachvault.Meta.Name)

			observer.StartVetoSlasherObserver(ctx, slasher, eachvault.Address) // in future updates this will change as if more type of slashable vaults are there
			// need to send data from veto slasher to alert parser

		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("🛑 Shutting down observers...")
	cancel()
	time.Sleep(2 * time.Second)
}
