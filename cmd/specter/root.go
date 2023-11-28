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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jinspect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "Input file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

		// Retrieve filepath of file supplied by global flag if exists
		file, _ := cmd.Flags().GetString("file")

		if file != "" {
			// Get the absolute Path of a file
			path, err := specter.GetAbsolutePath(file)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)

			} else {
				content, err := specter.ReadFromFile(path)

				if err != nil {
					fmt.Println("Somethin went wrong while reading the file.")
				}

				// This decision needs to either trigger the json or yaml parser later
				//specter.JsonOrYaml(content)

				fmt.Println(content)

				// more to happen here
			}

		} else {
			specter.ReadFromStdin()
		}
	},
}
