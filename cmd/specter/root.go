package specter

import (
	"fmt"
	"os"

	isp "github.com/betuxy/specter/internal/specter"
	"github.com/spf13/cobra"
)

var (
	file        string
	expanded    bool
	parsedData  []map[string]interface{}
	inputFormat string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "specter",
	Short: "Inspect JSON/YAML interactively",
	Long:  "A TUI tool to inspect JSON and YAML in an interactive tree view.\nReads from -f file or stdin if no file is given.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadInput(); err != nil {
			return err
		}
		return isp.RunTUI(parsedData, expanded)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Input file (reads stdin if omitted)")
	rootCmd.PersistentFlags().BoolVarP(&expanded, "expanded", "e", false, "Start with all nodes expanded")
	rootCmd.AddCommand(viewCmd, fmtCmd, convertCmd, versionCmd)
}

func loadInput() error {
	var (
		content []byte
		err     error
	)

	if file != "" {
		path, err := isp.GetAbsolutePath(file)
		if err != nil {
			return err
		}
		content, err = isp.ReadFromFile(path)
		if err != nil {
			return fmt.Errorf("can't read file: %w", err)
		}
	} else {
		content, err = isp.ReadFromStdin()
		if err != nil {
			return fmt.Errorf("can't read stdin: %w", err)
		}
	}

	if result, err := isp.ParseJSON(content); err == nil {
		parsedData = result
		inputFormat = "json"
	} else if result, err := isp.ParseYAML(content); err == nil {
		parsedData = result
		inputFormat = "yaml"
	} else {
		return fmt.Errorf("couldn't parse input as JSON or YAML")
	}

	return nil
}
