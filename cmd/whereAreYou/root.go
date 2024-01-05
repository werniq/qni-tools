package whereAreYou

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "IpTracker CLI application",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
