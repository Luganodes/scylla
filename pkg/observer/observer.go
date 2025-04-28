package observer

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/luganodes/slashing-observer/config"
	"github.com/luganodes/slashing-observer/pkg/alertmanager"
)

var abiJSON = `[
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "slashIndex", "type": "uint256" },
      { "indexed": false, "name": "slashedAmount", "type": "uint256" }
    ],
    "name": "ExecuteSlash",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "slashIndex", "type": "uint256" },
      { "indexed": true, "name": "subnetwork", "type": "bytes32" },
      { "indexed": true, "name": "operator", "type": "address" },
      { "indexed": false, "name": "slashAmount", "type": "uint256" },
      { "indexed": false, "name": "captureTimestamp", "type": "uint48" },
      { "indexed": false, "name": "vetoDeadline", "type": "uint48" }
    ],
    "name": "RequestSlash",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      { "indexed": true, "name": "slashIndex", "type": "uint256" },
      { "indexed": true, "name": "resolver", "type": "address" }
    ],
    "name": "VetoSlash",
    "type": "event"
  }
]`

func StartVetoSlasherObserver(ctx context.Context, address string) {
	client, err := ethclient.Dial(config.WS_URL)
	if err != nil {
		log.Printf("❌ [%s] Failed to connect to WebSocket: %v", address, err)
		return
	}

	contractAddr := common.HexToAddress(address)
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Printf("❌ [%s] Failed to parse ABI: %v", address, err)
		return
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics: [][]common.Hash{
			{
				parsedABI.Events["ExecuteSlash"].ID,
				parsedABI.Events["RequestSlash"].ID,
				parsedABI.Events["VetoSlash"].ID,
			},
		},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Printf("❌ [%s] Failed to subscribe to logs: %v", address, err)
		return
	}

	go func() {
		defer client.Close()
		log.Printf("📡 [%s] Listening for slashing events...", address)

		for {
			select {
			case err := <-sub.Err():
				log.Printf("❌ [%s] Subscription error: %v", address, err)
				return
			case vLog := <-logs:
				event, err := parsedABI.EventByID(vLog.Topics[0])
				if err != nil {
					log.Printf("❓ [%s] Unknown event: %s", address, vLog.Topics[0].Hex())
					continue
				}

				decoded := make(map[string]interface{})
				if err := parsedABI.UnpackIntoMap(decoded, event.Name, vLog.Data); err != nil {
					log.Printf("❌ [%s] Failed to unpack data: %v", address, err)
					continue
				}

				for i, input := range event.Inputs {
					if input.Indexed {
						decoded[input.Name] = parseTopic(vLog.Topics[i+1], input.Type.String())
					}
				}

				fmt.Printf("\n📣 [%s] Event: %s\n", address, event.Name) // from here to main  a channel
				for k, v := range decoded {
					fmt.Printf(" - %s: %v\n", k, v)
				}

				alertData := map[string]interface{}{
					"event":   event.Name,
					"address": address,
					"data":    decoded,
				}
				alertmanager.SendStructuredData(alertData)

			case <-ctx.Done():
				log.Printf("🛑 [%s] Stopping observer...", address)
				return
			}
		}
	}()
}

func parseTopic(topic common.Hash, typ string) interface{} {
	switch typ {
	case "address":
		return common.HexToAddress(topic.Hex())
	case "uint256":
		return new(big.Int).SetBytes(topic.Bytes())
	case "bytes32":
		return topic.Hex()
	default:
		return topic.Hex()
	}
}
