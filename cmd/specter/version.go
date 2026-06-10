package specter

import (
	"fmt"

	"github.com/spf13/cobra"
)

const appVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("specter %s\n", appVersion)
	},
}
