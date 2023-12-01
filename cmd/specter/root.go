package specter

import (
	"fmt"
	"os"

	"github.com/betuxy/specter/internal/specter"
	"github.com/spf13/cobra"
)

// Flags
var file string

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "Input file")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "specter",
	Short: "specter - inspect json objects interactively.",
	Long:  "A TUI tool to inspect JSON and Yaml objects in a tree view.",
	Run: func(cmd *cobra.Command, args []string) {

		// Store json or yaml object(s) in here
		var content []byte

		// Retrieve filepath of file supplied by global flag if exists
		file, _ := cmd.Flags().GetString("file")

		if file != "" {
			// Get the absolute Path of a file
			path, err := specter.GetAbsolutePath(file)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)

			} else {
				content, err = specter.ReadFromFile(path)

				if err != nil {
					fmt.Println("Error: Can't read from file.", err)
					os.Exit(1)
				}

			}

		} else {
			var err error
			content, err = specter.ReadFromStdin()
			if err != nil {
				fmt.Println("Error: Can't read from stdin.", err)
			}
		}

		var objectSlice []map[string]interface{}

		// Try to parse the input as json, then yaml, then print error and exit.
		if result, err := specter.ParseJSON(content); err == nil {
			objectSlice = result
		} else if result, err := specter.ParseYAML(content); err == nil {
			objectSlice = result
		} else {
			fmt.Println("Error: Wasn't able to parse either json or yaml")
			os.Exit(1)
		}

		specter.PrintAllJSON(objectSlice)


		// more to happen here

	},
}
