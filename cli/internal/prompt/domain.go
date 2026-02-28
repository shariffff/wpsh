package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wordmon/cli/internal/utils"
	"github.com/wordmon/cli/pkg/models"
)

// DomainAddInput holds the input for adding a domain
type DomainAddInput struct {
	ServerName string
	SiteID     string
	Domain     string
	IssueSSL   bool
}

// DomainRemoveInput holds the input for removing a domain
type DomainRemoveInput struct {
	ServerName string
	SiteID     string
	Domain     string
}

// DomainSSLInput holds the input for issuing SSL
type DomainSSLInput struct {
	ServerName   string
	SiteID       string
	Domain       string
	CertbotEmail string
}

// PromptDomainAdd prompts for domain addition details
func PromptDomainAdd(servers []models.Server) (*DomainAddInput, error) {
	input := &DomainAddInput{}

	// Build list of sites
	type SiteOption struct {
		ServerName string
		Site       models.Site
	}

	var siteOptions []SiteOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				siteOptions = append(siteOptions, SiteOption{
					ServerName: server.Name,
					Site:       site,
				})
			}
		}
	}

	if len(siteOptions) == 0 {
		return nil, fmt.Errorf("no sites available. Create a site first with: wordmon site create")
	}

	// Select site
	optionStrings := make([]string, len(siteOptions))
	for i, opt := range siteOptions {
		optionStrings[i] = fmt.Sprintf("%s on %s (%d domains)",
			opt.Site.PrimaryDomain, opt.ServerName, len(opt.Site.Domains))
	}

	var selectedIndex int
	selectPrompt := &survey.Select{
		Message: "Select site to add domain to:",
		Options: optionStrings,
		Help:    "Choose which WordPress site should serve this domain",
	}
	if err := survey.AskOne(selectPrompt, &selectedIndex); err != nil {
		return nil, err
	}

	input.ServerName = siteOptions[selectedIndex].ServerName
	input.SiteID = siteOptions[selectedIndex].Site.SiteID

	// Domain name
	domainPrompt := &survey.Input{
		Message: "Domain name to add:",
		Help:    "Enter the domain (e.g., www.example.com)",
	}
	if err := survey.AskOne(domainPrompt, &input.Domain, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidateDomain)); err != nil {
		return nil, err
	}

	// Check if domain already exists
	for _, domain := range siteOptions[selectedIndex].Site.Domains {
		if domain.Domain == input.Domain {
			return nil, fmt.Errorf("domain '%s' already exists on this site", input.Domain)
		}
	}

	// Ask about SSL
	sslPrompt := &survey.Confirm{
		Message: "Issue SSL certificate for this domain?",
		Default: true,
		Help:    "Automatically obtain a Let's Encrypt SSL certificate",
	}
	if err := survey.AskOne(sslPrompt, &input.IssueSSL); err != nil {
		return nil, err
	}

	return input, nil
}

// PromptDomainRemove prompts for domain removal
func PromptDomainRemove(servers []models.Server) (*DomainRemoveInput, error) {
	input := &DomainRemoveInput{}

	// Build list of all domains
	type DomainOption struct {
		ServerName string
		SiteID     string
		Domain     models.Domain
		IsPrimary  bool
	}

	var domainOptions []DomainOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				for _, domain := range site.Domains {
					domainOptions = append(domainOptions, DomainOption{
						ServerName: server.Name,
						SiteID:     site.SiteID,
						Domain:     domain,
						IsPrimary:  domain.Domain == site.PrimaryDomain,
					})
				}
			}
		}
	}

	if len(domainOptions) == 0 {
		return nil, fmt.Errorf("no domains available to remove")
	}

	// Create selection options
	optionStrings := make([]string, len(domainOptions))
	for i, opt := range domainOptions {
		sslStatus := ""
		if opt.Domain.SSLEnabled {
			sslStatus = " [SSL]"
		}
		primaryMarker := ""
		if opt.IsPrimary {
			primaryMarker = " (PRIMARY)"
		}
		optionStrings[i] = fmt.Sprintf("%s - %s on %s%s%s",
			opt.Domain.Domain, opt.SiteID, opt.ServerName, sslStatus, primaryMarker)
	}

	var selectedIndex int
	selectPrompt := &survey.Select{
		Message: "Select domain to remove:",
		Options: optionStrings,
		Help:    "Choose which domain to remove from the site",
	}
	if err := survey.AskOne(selectPrompt, &selectedIndex); err != nil {
		return nil, err
	}

	selected := domainOptions[selectedIndex]

	// Warn if removing primary domain
	if selected.IsPrimary {
		fmt.Println()
		fmt.Println("⚠️  WARNING: You are removing the PRIMARY domain for this site!")
		fmt.Println("This may break the WordPress installation.")
		fmt.Println()

		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Are you sure you want to remove the primary domain?",
			Default: false,
		}, &confirm); err != nil {
			return nil, err
		}

		if !confirm {
			return nil, fmt.Errorf("domain removal cancelled")
		}
	}

	input.ServerName = selected.ServerName
	input.SiteID = selected.SiteID
	input.Domain = selected.Domain.Domain

	return input, nil
}

// PromptDomainSSL prompts for SSL certificate issuance
func PromptDomainSSL(servers []models.Server, defaultEmail string) (*DomainSSLInput, error) {
	input := &DomainSSLInput{}

	// Build list of domains without SSL
	type DomainOption struct {
		ServerName string
		SiteID     string
		SiteDomain string
		Domain     models.Domain
	}

	var domainOptions []DomainOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				for _, domain := range site.Domains {
					if !domain.SSLEnabled {
						domainOptions = append(domainOptions, DomainOption{
							ServerName: server.Name,
							SiteID:     site.SiteID,
							SiteDomain: site.PrimaryDomain,
							Domain:     domain,
						})
					}
				}
			}
		}
	}

	if len(domainOptions) == 0 {
		return nil, fmt.Errorf("no domains without SSL certificates found")
	}

	// Create selection options
	optionStrings := make([]string, len(domainOptions))
	for i, opt := range domainOptions {
		optionStrings[i] = fmt.Sprintf("%s - site: %s on %s",
			opt.Domain.Domain, opt.SiteDomain, opt.ServerName)
	}

	var selectedIndex int
	selectPrompt := &survey.Select{
		Message: "Select domain to issue SSL for:",
		Options: optionStrings,
		Help:    "Choose which domain to obtain a Let's Encrypt certificate for",
	}
	if err := survey.AskOne(selectPrompt, &selectedIndex); err != nil {
		return nil, err
	}

	selected := domainOptions[selectedIndex]
	input.ServerName = selected.ServerName
	input.SiteID = selected.SiteID
	input.Domain = selected.Domain.Domain

	// Certbot email
	emailPrompt := &survey.Input{
		Message: "Email for Let's Encrypt notifications:",
		Default: defaultEmail,
		Help:    "Email address for certificate expiration notices",
	}
	if err := survey.AskOne(emailPrompt, &input.CertbotEmail, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidateEmail)); err != nil {
		return nil, err
	}

	return input, nil
}
