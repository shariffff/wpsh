package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wp-sh/cli/internal/config"
	"github.com/wp-sh/cli/internal/installer"
	"github.com/wp-sh/cli/internal/prompt"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize WPSH environment",
	Long: `Initialize WPSH by setting up configuration and copying Ansible playbooks.

This command will:
  1. Create ~/.wp-sh/ directory structure
  2. Copy Ansible playbooks from the repository to ~/.wp-sh/ansible/
  3. Create initial configuration file (wp-sh.yaml)
  4. Prompt for global settings (SSH key, certbot email)
  5. Validate the installation

Note: MySQL admin passwords are generated automatically per-server during provisioning.

Run this command once after installing WPSH.

Examples:
  # Interactive mode
  wp-sh init

  # Non-interactive mode
  wp-sh init --ssh-public-key ~/.ssh/id_rsa.pub --certbot-email admin@example.com

  # Force overwrite existing configuration
  wp-sh init --force`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  WPSH Initialization")
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		// Check if config already exists
		if mgr.ConfigExists() && !force {
			color.Yellow("Configuration file already exists at: %s", mgr.GetConfigPath())
			fmt.Println()
			fmt.Println("Options:")
			fmt.Printf("  • Edit the config:      %s %s\n", getEditor(), mgr.GetConfigPath())
			fmt.Println("  • Overwrite config:     wp-sh init --force")
			fmt.Println()
			fmt.Println("Use --force to overwrite the existing configuration.")
			os.Exit(1)
		}

		// Check if ansible is already initialized
		ansibleInitialized := installer.IsInitialized()

		if !ansibleInitialized {
			// Initialize ansible directory
			fmt.Print("→ Copying Ansible playbooks... ")
			if err := installer.Initialize(); err != nil {
				color.Red("✗")
				color.Red("\nError: %v", err)
				os.Exit(1)
			}
			color.Green("✓")
		} else {
			fmt.Println("→ Ansible playbooks already installed ✓")
		}

		// Get one-time setup values
		var initInput *prompt.InitInput

		// Check for non-interactive mode
		sshKey, _ := cmd.Flags().GetString("ssh-public-key")
		certbotEmail, _ := cmd.Flags().GetString("certbot-email")

		if sshKey != "" && certbotEmail != "" {
			// Non-interactive mode
			initInput = &prompt.InitInput{
				SSHPublicKey: sshKey,
				CertbotEmail: certbotEmail,
			}
		} else if sshKey != "" || certbotEmail != "" {
			// Partial flags provided
			color.Red("Error: All flags required for non-interactive mode: --ssh-public-key, --certbot-email")
			os.Exit(1)
		} else {
			// Interactive mode - prompt for setup values
			initInput, err = prompt.PromptInitSetup()
			if err != nil {
				color.Red("\nError: %v", err)
				os.Exit(1)
			}
		}

		// Create configuration with the provided values
		fmt.Print("→ Creating configuration file... ")

		cfg := config.DefaultConfig()
		cfg.Ansible.Path = installer.GetAnsibleDir()

		// Set global vars from user input
		cfg.GlobalVars["wp-sh_ssh_key"] = initInput.SSHPublicKey
		cfg.GlobalVars["certbot_email"] = initInput.CertbotEmail

		if err := mgr.Save(cfg); err != nil {
			color.Red("✗")
			color.Red("\nError: %v", err)
			os.Exit(1)
		}
		color.Green("✓")

		// Validate installation
		fmt.Print("→ Validating installation... ")
		if err := validateInstallation(); err != nil {
			color.Red("✗")
			color.Red("\nWarning: %v", err)
		} else {
			color.Green("✓")
		}

		// Success message
		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ WPSH initialized successfully!")
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Installation paths:")
		fmt.Printf("  • Ansible:       %s\n", installer.GetAnsibleDir())
		fmt.Printf("  • Config:        %s\n", mgr.GetConfigPath())
		fmt.Println()
		fmt.Println("Configuration saved:")
		fmt.Printf("  • SSH Public Key:    %s\n", initInput.SSHPublicKey)
		fmt.Printf("  • Certbot Email:     %s\n", initInput.CertbotEmail)
		fmt.Println()
		fmt.Println("To edit your configuration later:")
		fmt.Printf("  %s %s\n", getEditor(), mgr.GetConfigPath())
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Add a server:    wp-sh server add")
		fmt.Println("  2. Provision:       wp-sh server provision <name>")
		fmt.Println("  3. Create site:     wp-sh site create")
		fmt.Println()
	},
}

func validateInstallation() error {
	// Check if ansible directory has required files
	ansiblePath := installer.GetAnsibleDir()

	requiredFiles := []string{
		"provision.yml",
		"website.yml",
		"playbooks/domain_management.yml",
		"playbooks/delete_site.yml",
		"ansible.cfg",
	}

	for _, file := range requiredFiles {
		fullPath := fmt.Sprintf("%s/%s", ansiblePath, file)
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("missing required file: %s", file)
		}
	}

	return nil
}

// getEditor returns the user's preferred editor
func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "nano"
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags for non-interactive mode
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite existing configuration")
	initCmd.Flags().String("ssh-public-key", "", "Path to SSH public key for wp-sh user")
	initCmd.Flags().String("certbot-email", "", "Email for Let's Encrypt SSL certificates")
}
