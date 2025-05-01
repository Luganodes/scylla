package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/luganodes/slashing-observer/config"
	"github.com/luganodes/slashing-observer/pkg/observer"
	"github.com/luganodes/slashing-observer/pkg/schema"
	"github.com/luganodes/slashing-observer/pkg/vault"
)

type StartCmd struct{}

func (cmd *StartCmd) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	log.Println("🚀 Starting slashing observer...")

	// Track which vaults are already being observed
	observedVaults := make(map[string]struct{})
	var mu sync.Mutex

	// Function to fetch vaults and start new observers
	fetchAndStartObservers := func() {
		log.Printf("🔄 Fetching vaults...")
		vaults, err := vault.GetVaultInfoList()
		if err != nil {
			log.Println("❌ Failed to fetch vaults:", err)
			return
		}

		for _, eachvault := range vaults {
			if !eachvault.Slashable {
				continue
			}

			mu.Lock()
			if _, exists := observedVaults[eachvault.Address]; exists {
				mu.Unlock()
				continue
			}
			observedVaults[eachvault.Address] = struct{}{}
			mu.Unlock()

			go func(v schema.VaultInfo) {
				slasher, err := vault.GetSlasherAddress(config.RPC_URL, v.Address)
				if err != nil {
					log.Printf("❌ Error getting slasher address for vault %s: %v", v.Address, err)
					return
				}

				slasherType, err := vault.GetSlasherType(config.RPC_URL, slasher)
				if err != nil {
					log.Printf("❌ Error getting slasher type for vault %s: %v", v.Address, err)
					return
				}

				log.Printf("✅ New vault observer:\n  Vault: %s\n  Slasher: %s\n  Type: %d\n  Name: %s",
					v.Address, slasher, slasherType, v.Meta.Name)

				observer.StartVetoSlasherObserver(ctx, slasher)

			}(eachvault)
		}
	}

	// Initial run
	fetchAndStartObservers()

	// Start ticker for periodic fetching
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fetchAndStartObservers()
		case <-sig:
			log.Println("🛑 Shutting down observers...")
			cancel()
			time.Sleep(2 * time.Second) // allow goroutines to exit gracefully
			return nil
		}
	}
}
