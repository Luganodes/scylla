package alertmanager

import (
	"encoding/json"
	"fmt"
)

func SendStructuredData(data interface{}) {
	alertJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error converting data to JSON: %v\n", err)
		return
	}

	fmt.Printf("📣 ALERT DATA:\n%s\n", string(alertJSON))
}
