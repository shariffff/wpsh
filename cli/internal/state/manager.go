package state

import (
	"fmt"
	"time"

	"github.com/wordmon/cli/internal/config"
	"github.com/wordmon/cli/pkg/models"
)

// Manager handles state updates to the configuration
type Manager struct {
	configManager *config.Manager
}

// NewManager creates a new state manager
func NewManager(configManager *config.Manager) *Manager {
	return &Manager{
		configManager: configManager,
	}
}

// MarkServerProvisioned updates a server's status to provisioned
func (m *Manager) MarkServerProvisioned(serverName string) error {
	// Load current config
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find and update server
	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			now := time.Now()
			cfg.Servers[i].Status = "provisioned"
			cfg.Servers[i].ProvisionedAt = &now
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("server not found: %s", serverName)
	}

	// Save updated config
	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// MarkServerError updates a server's status to error
func (m *Manager) MarkServerError(serverName string) error {
	// Load current config
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find and update server
	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			cfg.Servers[i].Status = "error"
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("server not found: %s", serverName)
	}

	// Save updated config
	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// GetServer retrieves a server by name
func (m *Manager) GetServer(serverName string) (*models.Server, error) {
	cfg, err := m.configManager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	for _, server := range cfg.Servers {
		if server.Name == serverName {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("server not found: %s", serverName)
}

// AddSiteToServer adds a site to a server's configuration
func (m *Manager) AddSiteToServer(serverName string, site models.Site) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			cfg.Servers[i].Sites = append(cfg.Servers[i].Sites, site)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("server not found: %s", serverName)
	}

	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RemoveSiteFromServer removes a site from a server's configuration
func (m *Manager) RemoveSiteFromServer(serverName string, siteID string) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			// Filter out the site to remove
			newSites := make([]models.Site, 0)
			for _, site := range cfg.Servers[i].Sites {
				if site.SiteID != siteID {
					newSites = append(newSites, site)
				} else {
					found = true
				}
			}
			cfg.Servers[i].Sites = newSites
			break
		}
	}

	if !found {
		return fmt.Errorf("site '%s' not found on server '%s'", siteID, serverName)
	}

	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// AddDomainToSite adds a domain to a site's configuration
func (m *Manager) AddDomainToSite(serverName string, siteID string, domain models.Domain) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			for j := range cfg.Servers[i].Sites {
				if cfg.Servers[i].Sites[j].SiteID == siteID {
					cfg.Servers[i].Sites[j].Domains = append(cfg.Servers[i].Sites[j].Domains, domain)
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("site '%s' not found on server '%s'", siteID, serverName)
	}

	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RemoveDomainFromSite removes a domain from a site's configuration
func (m *Manager) RemoveDomainFromSite(serverName string, siteID string, domainName string) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			for j := range cfg.Servers[i].Sites {
				if cfg.Servers[i].Sites[j].SiteID == siteID {
					// Filter out the domain to remove
					newDomains := make([]models.Domain, 0)
					for _, d := range cfg.Servers[i].Sites[j].Domains {
						if d.Domain != domainName {
							newDomains = append(newDomains, d)
						} else {
							found = true
						}
					}
					cfg.Servers[i].Sites[j].Domains = newDomains
					break
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("domain '%s' not found on site '%s' on server '%s'", domainName, siteID, serverName)
	}

	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// UpdateDomainSSL updates a domain's SSL information
func (m *Manager) UpdateDomainSSL(serverName string, siteID string, domainName string, updatedDomain models.Domain) error {
	cfg, err := m.configManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	found := false
	for i := range cfg.Servers {
		if cfg.Servers[i].Name == serverName {
			for j := range cfg.Servers[i].Sites {
				if cfg.Servers[i].Sites[j].SiteID == siteID {
					for k := range cfg.Servers[i].Sites[j].Domains {
						if cfg.Servers[i].Sites[j].Domains[k].Domain == domainName {
							cfg.Servers[i].Sites[j].Domains[k] = updatedDomain
							found = true
							break
						}
					}
					break
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("domain '%s' not found on site '%s' on server '%s'", domainName, siteID, serverName)
	}

	if err := m.configManager.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
