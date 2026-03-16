package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wp-sh/cli/internal/ansible"
	"github.com/wp-sh/cli/internal/config"
	"github.com/wp-sh/cli/internal/prompt"
	"github.com/wp-sh/cli/internal/state"
	"github.com/wp-sh/cli/internal/utils"
	"github.com/wp-sh/cli/pkg/models"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains and SSL certificates",
	Long:  `Add domains to sites, remove domains, and issue SSL certificates.`,
}

// domainAddCmd represents the domain add command
var domainAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "Add a domain to a site",
	Long: `Add a new domain to an existing WordPress site and optionally issue an SSL certificate.

Examples:
  # Interactive mode
  wp-sh domain add

  # Non-interactive mode (for automation/AI agents)
  wp-sh domain add --server myserver --site mysite --domain www.example.com --ssl`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			outputError(cmd, "Failed to create config manager", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			outputError(cmd, "Configuration file not found", fmt.Errorf("run 'wp-sh init' first"))
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			outputError(cmd, "Failed to load configuration", err)
			os.Exit(1)
		}

		var input *prompt.DomainAddInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			issueSSL, _ := cmd.Flags().GetBool("ssl")
			input = &prompt.DomainAddInput{
				ServerName: serverName,
				SiteID:     siteName,
				Domain:     domain,
				IssueSSL:   issueSSL,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainAdd(cfg.Servers)
			if err != nil {
				outputError(cmd, "Failed to get domain details", err)
				os.Exit(1)
			}
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"operation": "add_domain",
			"domain":    input.Domain,
			"site_id":   input.SiteID,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Adding domain: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Domain addition failed: %v", err)
			os.Exit(1)
		}

		// Add domain to configuration
		newDomain := models.Domain{
			Domain:     input.Domain,
			SSLEnabled: false,
		}

		stateMgr := state.NewManager(mgr)
		if err := stateMgr.AddDomainToSite(input.ServerName, input.SiteID, newDomain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		color.Green("\n✓ Domain '%s' added successfully", input.Domain)

		// Issue SSL if requested
		if input.IssueSSL {
			fmt.Println()
			color.Cyan("═══════════════════════════════════════════════════════")
			color.Cyan("  Issuing SSL certificate for: %s", input.Domain)
			color.Cyan("═══════════════════════════════════════════════════════")
			fmt.Println()

			// Get certbot email from global vars
			certbotEmail := "admin@example.com"
			if email, ok := cfg.GlobalVars["certbot_email"].(string); ok {
				certbotEmail = email
			}

			sslVars := map[string]interface{}{
				"operation":     "issue_ssl",
				"domain":        input.Domain,
				"certbot_email": certbotEmail,
			}

			sslResult, err := executor.ExecutePlaybookWithResult("playbooks/domain_management.yml", *targetServer, sslVars, cfg.GlobalVars)
			if err != nil {
				color.Red("\n✗ SSL certificate issuance failed: %v", err)
				fmt.Println("The domain has been added but SSL is not configured.")
				fmt.Println("You can issue SSL later with: wp-sh domain ssl")
				os.Exit(1)
			}

			// Update domain with SSL info
			now := time.Now()
			var expiresAt *time.Time

			// Try to parse actual expiry from Ansible output
			if sslResult.SSLInfo != nil && sslResult.SSLInfo.Expiry != "" {
				expiresAt = utils.ParseSSLExpiry(sslResult.SSLInfo.Expiry)
			}

			// Fallback to 90 days if parsing fails
			if expiresAt == nil {
				fallback := now.AddDate(0, 3, 0)
				expiresAt = &fallback
			}

			sslDomain := models.Domain{
				Domain:       input.Domain,
				SSLEnabled:   true,
				SSLIssuedAt:  &now,
				SSLExpiresAt: expiresAt,
			}

			if err := stateMgr.UpdateDomainSSL(input.ServerName, input.SiteID, input.Domain, sslDomain); err != nil {
				color.Red("Warning: Failed to update SSL status in configuration: %v", err)
			}

			color.Green("\n✓ SSL certificate issued successfully")
			fmt.Println()
			fmt.Printf("Domain URL:  https://%s\n", input.Domain)
			fmt.Printf("Expires:     %s\n", expiresAt.Format("2006-01-02"))
		} else {
			fmt.Println()
			fmt.Printf("Domain URL:  http://%s\n", input.Domain)
			fmt.Println()
			fmt.Println("To issue SSL later: wp-sh domain ssl")
		}
	},
}

// domainRemoveCmd represents the domain remove command
var domainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete"},
	Short:   "Remove a domain from a site",
	Long: `Remove a domain from a WordPress site and its Nginx configuration.

Examples:
  # Interactive mode
  wp-sh domain remove

  # Non-interactive mode (for automation/AI agents)
  wp-sh domain remove --server myserver --site mysite --domain www.example.com --force`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			outputError(cmd, "Failed to create config manager", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			outputError(cmd, "Configuration file not found", fmt.Errorf("run 'wp-sh init' first"))
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			outputError(cmd, "Failed to load configuration", err)
			os.Exit(1)
		}

		var input *prompt.DomainRemoveInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			input = &prompt.DomainRemoveInput{
				ServerName: serverName,
				SiteID:     siteName,
				Domain:     domain,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainRemove(cfg.Servers)
			if err != nil {
				outputError(cmd, "Failed to get domain details", err)
				os.Exit(1)
			}
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Final confirmation
		color.Yellow("\n⚠️  WARNING: This will remove:")
		fmt.Printf("  - Domain: %s\n", input.Domain)
		fmt.Printf("  - Nginx configuration\n")
		fmt.Printf("  - SSL certificate (if any)\n")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Remove this domain?",
				Default: false,
			}, &confirm); err != nil {
				os.Exit(1)
			}

			if !confirm {
				fmt.Println("Domain removal cancelled")
				return
			}
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"operation": "remove_domain",
			"domain":    input.Domain,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Removing domain: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Domain removal failed: %v", err)
			os.Exit(1)
		}

		// Remove domain from configuration
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.RemoveDomainFromSite(input.ServerName, input.SiteID, input.Domain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		color.Green("\n✓ Domain '%s' removed successfully", input.Domain)
	},
}

// domainSSLCmd represents the domain ssl command
var domainSSLCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Issue SSL certificate for a domain",
	Long: `Obtain a Let's Encrypt SSL certificate for a domain.

Examples:
  # Interactive mode
  wp-sh domain ssl

  # Non-interactive mode (for automation/AI agents)
  wp-sh domain ssl --server myserver --site mysite --domain www.example.com --email admin@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			outputError(cmd, "Failed to create config manager", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			outputError(cmd, "Configuration file not found", fmt.Errorf("run 'wp-sh init' first"))
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			outputError(cmd, "Failed to load configuration", err)
			os.Exit(1)
		}

		// Get default certbot email from config
		defaultEmail := "admin@example.com"
		if email, ok := cfg.GlobalVars["certbot_email"].(string); ok {
			defaultEmail = email
		}

		var input *prompt.DomainSSLInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			email, _ := cmd.Flags().GetString("email")
			if email == "" {
				email = defaultEmail
			}
			input = &prompt.DomainSSLInput{
				ServerName:   serverName,
				SiteID:       siteName,
				Domain:       domain,
				CertbotEmail: email,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainSSL(cfg.Servers, defaultEmail)
			if err != nil {
				outputError(cmd, "Failed to get SSL details", err)
				os.Exit(1)
			}
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"operation":     "issue_ssl",
			"domain":        input.Domain,
			"certbot_email": input.CertbotEmail,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Issuing SSL certificate for: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		result, err := executor.ExecutePlaybookWithResult("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars)
		if err != nil {
			color.Red("\n✗ SSL certificate issuance failed: %v", err)
			os.Exit(1)
		}

		// Update domain with SSL info
		now := time.Now()
		var expiresAt *time.Time

		// Try to parse actual expiry from Ansible output
		if result.SSLInfo != nil && result.SSLInfo.Expiry != "" {
			expiresAt = utils.ParseSSLExpiry(result.SSLInfo.Expiry)
		}

		// Fallback to 90 days if parsing fails
		if expiresAt == nil {
			fallback := now.AddDate(0, 3, 0)
			expiresAt = &fallback
		}

		sslDomain := models.Domain{
			Domain:       input.Domain,
			SSLEnabled:   true,
			SSLIssuedAt:  &now,
			SSLExpiresAt: expiresAt,
		}

		stateMgr := state.NewManager(mgr)
		if err := stateMgr.UpdateDomainSSL(input.ServerName, input.SiteID, input.Domain, sslDomain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ SSL certificate issued successfully!")
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Printf("Domain:      https://%s\n", input.Domain)
		fmt.Printf("Issued:      %s\n", now.Format("2006-01-02"))
		fmt.Printf("Expires:     %s\n", expiresAt.Format("2006-01-02"))
		fmt.Printf("Auto-renew:  Certbot will auto-renew before expiration\n")
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainRemoveCmd)
	domainCmd.AddCommand(domainSSLCmd)

	// domain add flags (non-interactive mode)
	domainAddCmd.Flags().String("server", "", "Server name")
	domainAddCmd.Flags().String("site", "", "Site ID")
	domainAddCmd.Flags().String("domain", "", "Domain to add")
	domainAddCmd.Flags().Bool("ssl", false, "Issue SSL certificate for the domain")
	domainAddCmd.Flags().Bool("json", false, "Output in JSON format")

	// domain remove flags
	domainRemoveCmd.Flags().String("server", "", "Server name")
	domainRemoveCmd.Flags().String("site", "", "Site ID")
	domainRemoveCmd.Flags().String("domain", "", "Domain to remove")
	domainRemoveCmd.Flags().BoolP("force", "f", false, "Force removal without confirmation")
	domainRemoveCmd.Flags().Bool("json", false, "Output in JSON format")

	// domain ssl flags (non-interactive mode)
	domainSSLCmd.Flags().String("server", "", "Server name")
	domainSSLCmd.Flags().String("site", "", "Site ID")
	domainSSLCmd.Flags().String("domain", "", "Domain to issue SSL for")
	domainSSLCmd.Flags().String("email", "", "Email for Let's Encrypt notifications")
	domainSSLCmd.Flags().Bool("json", false, "Output in JSON format")
}
