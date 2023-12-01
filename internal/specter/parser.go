package specter

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseJSON(data []byte) ([]map[string]interface{}, error) {
	var jsonArray []map[string]interface{}
	var jsonObject map[string]interface{}

	if err := json.Unmarshal(data, &jsonArray); err == nil {
		return jsonArray, nil
	}

	if err := json.Unmarshal(data, &jsonObject); err != nil {
		return nil, err
	}

	return []map[string]interface{}{jsonObject}, nil
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
