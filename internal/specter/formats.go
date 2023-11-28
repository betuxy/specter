package specter

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// This is not finished at all

func JsonOrYaml(content []byte) bool {
	isJSON := isJSONFormat(content)
	isYAML := isYAMLFormat(content)

	if isJSON {
		return true
	} else if isYAML {
		return false
	}
	return false
}

func isJSONFormat(data []byte) bool {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	return err == nil
}

func isYAMLFormat(data []byte) bool {
	var yamlData map[string]interface{}
	err := yaml.Unmarshal([]byte(data), &yamlData)
	return err == nil
}
