package specter

import (
	isp "github.com/betuxy/specter/internal/specter"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Open an interactive TUI tree view (default)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadInput(); err != nil {
			return err
		}
		return isp.RunTUI(parsedData, expanded)
	},
}
