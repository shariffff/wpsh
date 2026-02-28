package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"

	// Global flags
	Verbose bool
	DryRun  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "wordmon",
	Short: "WordMon - Ansible wrapper for WordPress hosting management",
	Long: `WordMon is a CLI tool that simplifies WordPress hosting management
by wrapping Ansible playbooks with an intuitive, interactive interface.

Manage servers, sites, and domains with ease while maintaining full
visibility into your infrastructure state via ~/.wordmon/wordmon.yaml

Examples:
  # Initialize configuration
  wordmon init

  # Add and provision a new server
  wordmon server provision

  # Create a WordPress site
  wordmon site create

  # Add a domain with SSL
  wordmon domain add

  # List all servers
  wordmon server list --json`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "Show what would be done without making changes")
}
