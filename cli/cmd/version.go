package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version, commit SHA, and build date for wordmon CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("wordmon version %s\n", Version)
		fmt.Printf("Commit: %s\n", CommitSHA)
		fmt.Printf("Built: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
