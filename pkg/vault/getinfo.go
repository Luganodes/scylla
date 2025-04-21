package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/luganodes/slashing-observer/config"
	"github.com/luganodes/slashing-observer/pkg/schema"
)

func GetVaultInfoList() ([]schema.VaultInfo, error) {

	resp, err := http.Get(config.API_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var vaults []schema.VaultInfo
	if err := json.Unmarshal(body, &vaults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return vaults, nil
}
