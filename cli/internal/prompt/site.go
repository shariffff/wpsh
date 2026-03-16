package prompt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wpsh/cli/internal/utils"
	"github.com/wpsh/cli/pkg/models"
)

// SiteInput holds the input for site creation
type SiteInput struct {
	ServerName    string
	Domain        string
	SiteID        string
	AdminUser     string
	AdminEmail    string
	AdminPassword string
}

// PromptSiteCreate prompts for site creation details
func PromptSiteCreate(servers []models.Server) (*SiteInput, error) {
	input := &SiteInput{}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available. Add a server first with: wpsh server add")
	}

	// Filter only provisioned servers
	provisionedServers := make([]models.Server, 0)
	for _, s := range servers {
		if s.Status == "provisioned" {
			provisionedServers = append(provisionedServers, s)
		}
	}

	if len(provisionedServers) == 0 {
		return nil, fmt.Errorf("no provisioned servers available. Provision a server first with: wpsh server provision <name>")
	}

	// 1. Select server
	serverOptions := make([]string, len(provisionedServers))
	for i, s := range provisionedServers {
		serverOptions[i] = fmt.Sprintf("%s (%s) - %d sites", s.Name, s.IP, len(s.Sites))
	}

	var serverIndex int
	serverPrompt := &survey.Select{
		Message: "Select target server:",
		Options: serverOptions,
		Help:    "Choose a provisioned server to host this WordPress site",
	}
	if err := survey.AskOne(serverPrompt, &serverIndex); err != nil {
		return nil, err
	}
	input.ServerName = provisionedServers[serverIndex].Name

	// 2. Domain name
	domainPrompt := &survey.Input{
		Message: "Primary domain name:",
		Help:    "The main domain for this WordPress site (e.g., example.com)",
	}
	if err := survey.AskOne(domainPrompt, &input.Domain, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidateDomain)); err != nil {
		return nil, err
	}

	// 3. Auto-generate unique site ID (no prompt needed)
	selectedServer := provisionedServers[serverIndex]
	input.SiteID = generateUniqueSiteID(input.Domain, selectedServer.Sites)

	// 4. WordPress admin user
	adminUserPrompt := &survey.Input{
		Message: "WordPress admin username:",
		Default: "admin",
		Help:    "Username for WordPress admin account",
	}
	if err := survey.AskOne(adminUserPrompt, &input.AdminUser, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// 5. WordPress admin email
	adminEmailPrompt := &survey.Input{
		Message: "WordPress admin email:",
		Help:    "Email address for WordPress admin account",
	}
	if err := survey.AskOne(adminEmailPrompt, &input.AdminEmail, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidateEmail)); err != nil {
		return nil, err
	}

	// 6. WordPress admin password (with option to generate)
	var useGeneratedPassword bool
	generatePrompt := &survey.Confirm{
		Message: "Generate secure password?",
		Default: true,
		Help:    "Auto-generate a strong password or enter your own",
	}
	if err := survey.AskOne(generatePrompt, &useGeneratedPassword); err != nil {
		return nil, err
	}

	if useGeneratedPassword {
		input.AdminPassword = GenerateSecurePassword(20)
		fmt.Printf("\n")
		fmt.Printf("Generated password: %s\n", input.AdminPassword)
		fmt.Printf("⚠️  IMPORTANT: Save this password securely!\n")
		fmt.Printf("\n")

		var acknowledged bool
		ackPrompt := &survey.Confirm{
			Message: "Have you saved the password?",
			Default: false,
		}
		if err := survey.AskOne(ackPrompt, &acknowledged); err != nil {
			return nil, err
		}
		if !acknowledged {
			return nil, fmt.Errorf("please save the password before continuing")
		}
	} else {
		passwordPrompt := &survey.Password{
			Message: "WordPress admin password:",
			Help:    "Min 12 chars with uppercase, lowercase, number, and special character",
		}
		if err := survey.AskOne(passwordPrompt, &input.AdminPassword, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidatePasswordStrength)); err != nil {
			return nil, err
		}
	}

	// 7. Confirmation
	if err := confirmSiteCreation(input); err != nil {
		return nil, err
	}

	return input, nil
}

// confirmSiteCreation shows a summary and asks for confirmation
func confirmSiteCreation(input *SiteInput) error {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("Site Configuration Summary:")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("  Server:       %s\n", input.ServerName)
	fmt.Printf("  Domain:       %s\n", input.Domain)
	fmt.Printf("  Site ID:      %s\n", input.SiteID)
	fmt.Printf("  Admin User:   %s\n", input.AdminUser)
	fmt.Printf("  Admin Email:  %s\n", input.AdminEmail)
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Create this WordPress site?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("site creation cancelled")
	}

	return nil
}

// Helper functions

// generateUniqueSiteID creates a unique site ID from the domain, handling collisions
func generateUniqueSiteID(domain string, existingSites []models.Site) string {
	base := generateBaseSiteID(domain)

	// Check for collisions and append number if needed
	candidate := base
	suffix := 2
	for siteIDExists(existingSites, candidate) {
		// Calculate how much space we have for the suffix
		suffixStr := fmt.Sprintf("%d", suffix)
		maxBaseLen := 16 - len(suffixStr)
		if len(base) > maxBaseLen {
			candidate = base[:maxBaseLen] + suffixStr
		} else {
			candidate = base + suffixStr
		}
		suffix++
	}

	return candidate
}

// generateBaseSiteID generates a base site ID from a domain name
func generateBaseSiteID(domain string) string {
	// Remove common TLDs
	tlds := []string{".com", ".net", ".org", ".io", ".co", ".dev", ".app", ".xyz", ".info", ".biz"}
	name := domain
	for _, tld := range tlds {
		name = strings.TrimSuffix(name, tld)
	}

	// Remove all non-alphanumeric characters
	alphanumRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	name = alphanumRegex.ReplaceAllString(name, "")

	// Limit to 14 characters (leaving room for numeric suffix)
	if len(name) > 14 {
		name = name[:14]
	}

	// Ensure at least 3 characters
	if len(name) < 3 {
		name = "site" + name
	}

	return strings.ToLower(name)
}

// siteIDExists checks if a site ID already exists in the list
func siteIDExists(sites []models.Site, id string) bool {
	for _, site := range sites {
		if site.SiteID == id {
			return true
		}
	}
	return false
}

// GenerateSiteID is exported for use by cmd package in non-interactive mode
func GenerateSiteID(domain string, existingSites []models.Site) string {
	return generateUniqueSiteID(domain, existingSites)
}

// GenerateSecurePassword generates a cryptographically secure random password
func GenerateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a simple method if crypto/rand fails
			password[i] = charset[i%len(charset)]
		} else {
			password[i] = charset[num.Int64()]
		}
	}

	return string(password)
}
