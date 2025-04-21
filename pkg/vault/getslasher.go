package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const vaultABIJSON = `
[
  {
    "constant": true,
    "inputs": [],
    "name": "slasher",
    "outputs": [
      {
        "name": "",
        "type": "address"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  }
]`

const slasherABIJSON = `
[
  {
    "constant": true,
    "inputs": [],
    "name": "TYPE",
    "outputs": [
      {
        "name": "",
        "type": "uint64"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  }
]`

var implSlot = common.HexToHash("0x360894A13BA1A3210667C828492DB98DCA3E2076CC3735A920A3CA505D382BBC")

func GetSlasherAddress(rpcURL, proxyAddress string) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	defer client.Close()

	proxy := common.HexToAddress(proxyAddress)

	ctx := context.Background()
	storage, err := client.StorageAt(ctx, proxy, implSlot, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get implementation storage: %w", err)
	}

	// last 20 bytes is the address
	if len(storage) < 32 {
		return "", fmt.Errorf("unexpected storage size")
	}
	_ = common.BytesToAddress(storage[12:])

	// Parse ABI and pack the call
	parsedABI, err := abi.JSON(strings.NewReader(vaultABIJSON))
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("slasher")
	if err != nil {
		return "", fmt.Errorf("failed to pack method: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &proxy,
		Data: data,
	}
	res, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("contract call failed: %w", err)
	}

	var slasher common.Address
	if err := parsedABI.UnpackIntoInterface(&slasher, "slasher", res); err != nil {
		return "", fmt.Errorf("failed to unpack result: %w", err)
	}

	return slasher.Hex(), nil
}

func GetSlasherType(rpcURL, slasherAddress string) (uint64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	defer client.Close()

	address := common.HexToAddress(slasherAddress)
	ctx := context.Background()

	parsedABI, err := abi.JSON(strings.NewReader(slasherABIJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to parse Slasher ABI: %w", err)
	}

	data, err := parsedABI.Pack("TYPE")
	if err != nil {
		return 0, fmt.Errorf("failed to pack TYPE call: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

	res, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return 0, fmt.Errorf("contract call to slasher failed: %w", err)
	}

	var result uint64
	if err := parsedABI.UnpackIntoInterface(&result, "TYPE", res); err != nil {
		return 0, fmt.Errorf("failed to unpack TYPE result: %w", err)
	}

	return result, nil
}
