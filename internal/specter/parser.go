package specter

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// This is not finished at all
func JsonOrYaml(content []byte) (map[string]interface{}, error) {
	json := make(map[string]interface{})
	yaml := make(map[string]interface{})
	if _, jsonErr := isJSONFormat(content, json); jsonErr == nil {
		return json, nil
	} else if _, yamlErr := isYAMLFormat(content, yaml); yamlErr == nil {
		return yaml, nil
	} else {
		return nil, jsonErr
	}
}

func isJSONFormat(data []byte, jsonData *map[string]interface{}) (bool, error) {
	if err := json.Unmarshal(data, jsonData); err != nil {
		return false, err
	}
	return true, nil
}

func isYAMLFormat(data []byte, yamlData map[string]interface{}) (bool, error) {
	err := yaml.Unmarshal(data, yamlData)
	return false, err
}

func PrintSerializedJson(data map[string]interface{}) {
	// Serialize map data back into JSON
	serializedJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}

	fmt.Println(string(serializedJSON))
}
