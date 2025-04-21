package schema

type VaultInfo struct {
	Address   string `json:"address"`
	Slashable bool   `json:"slashable"`
	Meta      struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"meta"`
}
