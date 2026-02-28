package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wordmon/cli/internal/config"
	"github.com/wordmon/cli/internal/prompt"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage wordmon configuration",
	Long:  `Display, validate, and edit the wordmon configuration file.`,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display the contents of the wordmon configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found at: %s", mgr.GetConfigPath())
			fmt.Println("Run 'wordmon init' to create it.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Marshal to YAML for pretty display
		data, err := yaml.Marshal(cfg)
		if err != nil {
			color.Red("Error: Failed to marshal configuration: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration file: %s\n\n", mgr.GetConfigPath())
		fmt.Println(string(data))
	},
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Validate the wordmon configuration file for correctness and consistency.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found at: %s", mgr.GetConfigPath())
			fmt.Println("Run 'wordmon init' to create it.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		validator := config.NewValidator()

		// Validate struct
		fmt.Println("Validating configuration structure...")
		if err := validator.ValidateStruct(cfg); err != nil {
			color.Red("✗ Structure validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Structure validation passed")

		// Validate business rules
		fmt.Println("Validating business rules...")
		if err := validator.ValidateBusinessRules(cfg); err != nil {
			color.Red("✗ Business rules validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Business rules validation passed")

		// Validate Ansible environment
		fmt.Println("Validating Ansible environment...")
		if err := validator.ValidateAnsibleEnvironment(cfg); err != nil {
			color.Red("✗ Ansible environment validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Ansible environment validation passed")

		fmt.Println()
		color.Green("✓ Configuration is valid")
		fmt.Printf("  Servers: %d\n", len(cfg.Servers))

		totalSites := 0
		for _, server := range cfg.Servers {
			totalSites += len(server.Sites)
		}
		fmt.Printf("  Sites: %d\n", totalSites)
	},
}

// configEditCmd represents the config edit command
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration file in your preferred editor",
	Long:  `Open the wordmon configuration file in your preferred editor. On first run, you'll be prompted to select an editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found at: %s", mgr.GetConfigPath())
			fmt.Println("Run 'wordmon init' to create it.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// If no preferred editor is set, prompt for one
		if cfg.PreferredEditor == "" {
			editor, err := prompt.PromptEditorSelection()
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}

			cfg.PreferredEditor = editor
			if err := mgr.Save(cfg); err != nil {
				color.Red("Error: Failed to save editor preference: %v", err)
				os.Exit(1)
			}
			color.Green("Saved editor preference: %s", editor)
		}

		// Open the config file in the preferred editor
		editorCmd := exec.Command(cfg.PreferredEditor, mgr.GetConfigPath())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			color.Red("Error: Failed to open editor: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configEditCmd)
}
