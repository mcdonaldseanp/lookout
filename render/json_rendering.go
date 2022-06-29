package render

import (
	"encoding/json"
	"fmt"
)

func RenderJson(data interface{}) (string, error) {
	json_output, json_err := json.Marshal(data)
	if json_err != nil {
		return "", fmt.Errorf("Could not render result as JSON: %s\n", json_err)
	}
	return string(json_output), nil
}
