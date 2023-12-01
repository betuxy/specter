package specter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

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

func ParseYAML(data []byte) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))

	for {
		var document map[string]interface{}
		if err := decoder.Decode(&document); err != nil {
			if err == io.EOF {
				break // End of documents
			}
			return nil, err
		}

		result = append(result, document)
	}

	return result, nil
}

func PrintAllJSON(data []map[string]interface{}) {
	objects, err := toJSON(data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, str := range objects {
		fmt.Println(str)
	}
}

func toJSON(data []map[string]interface{}) ([]string, error) {
	var objects []string
	for _, object := range data {
		serializedJSON, err := json.Marshal(object)
		if err != nil {
			continue
		}
		objects = append(objects, string(serializedJSON))
	}

	return objects, nil
}
