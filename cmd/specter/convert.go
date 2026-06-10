package specter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var targetFormat string

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between JSON and YAML",
	Example: "  specter convert -f file.json --to yaml\n  cat file.yaml | specter convert --to json",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadInput(); err != nil {
			return err
		}
		return convertOutput()
	},
}

func init() {
	convertCmd.Flags().StringVarP(&targetFormat, "to", "t", "", "Target format: json or yaml (required)")
	convertCmd.MarkFlagRequired("to")
}

func convertOutput() error {
	switch strings.ToLower(targetFormat) {
	case "yaml":
		for i, obj := range parsedData {
			if i > 0 {
				fmt.Println("---")
			}
			out, err := yaml.Marshal(obj)
			if err != nil {
				return err
			}
			fmt.Print(string(out))
		}
	case "json":
		if len(parsedData) == 1 {
			out, err := json.MarshalIndent(parsedData[0], "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		} else {
			out, err := json.MarshalIndent(parsedData, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		}
	default:
		return fmt.Errorf("unknown format %q: must be json or yaml", targetFormat)
	}
	return nil
}
