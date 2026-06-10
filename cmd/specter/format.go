package specter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var indentSize int

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Pretty-print JSON or YAML to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadInput(); err != nil {
			return err
		}
		return prettyPrint()
	},
}

func init() {
	fmtCmd.Flags().IntVarP(&indentSize, "indent", "i", 2, "Number of spaces for JSON indentation")
}

func prettyPrint() error {
	indent := strings.Repeat(" ", indentSize)
	for i, obj := range parsedData {
		switch inputFormat {
		case "yaml":
			if i > 0 {
				fmt.Println("---")
			}
			out, err := yaml.Marshal(obj)
			if err != nil {
				return err
			}
			fmt.Print(string(out))
		default:
			out, err := json.MarshalIndent(obj, "", indent)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		}
	}
	return nil
}
