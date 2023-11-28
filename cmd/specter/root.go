package specter

import (
	"os"

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
	Use:   "jinspect",
	Short: "Inspect json objects interactively.",
	Long:  "A TUI tool to inspect JSON and Yaml objects in a tree view.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
