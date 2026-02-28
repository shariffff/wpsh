package utils

import (
	"time"

	"github.com/wordmon/cli/pkg/models"
)

// FindServerByName finds a server by name in the servers slice
// Returns nil if not found
func FindServerByName(servers []models.Server, name string) *models.Server {
	for i := range servers {
		if servers[i].Name == name {
			return &servers[i]
		}
	}
	return nil
}

// FindServerIndexByName finds a server's index by name in the servers slice
// Returns -1 if not found
func FindServerIndexByName(servers []models.Server, name string) int {
	for i := range servers {
		if servers[i].Name == name {
			return i
		}
	}
	return -1
}

// FindSiteBySiteID finds a site by site ID within a server
// Returns nil if not found
func FindSiteBySiteID(server *models.Server, siteID string) *models.Site {
	if server == nil {
		return nil
	}
	for i := range server.Sites {
		if server.Sites[i].SiteID == siteID {
			return &server.Sites[i]
		}
	}
	return nil
}

// FindSiteIndexBySiteID finds a site's index by site ID within a server
// Returns -1 if not found
func FindSiteIndexBySiteID(server *models.Server, siteID string) int {
	if server == nil {
		return -1
	}
	for i := range server.Sites {
		if server.Sites[i].SiteID == siteID {
			return i
		}
	}
	return -1
}

// FindSiteByDomain finds a site by domain (primary or additional) within a server
// Returns nil if not found
func FindSiteByDomain(server *models.Server, domain string) *models.Site {
	if server == nil {
		return nil
	}
	for i := range server.Sites {
		if server.Sites[i].PrimaryDomain == domain {
			return &server.Sites[i]
		}
		for _, d := range server.Sites[i].Domains {
			if d.Domain == domain {
				return &server.Sites[i]
			}
		}
	}
	return nil
}

// GetProvisionedServers returns only servers with status "provisioned"
func GetProvisionedServers(servers []models.Server) []models.Server {
	result := make([]models.Server, 0)
	for _, s := range servers {
		if s.Status == "provisioned" {
			result = append(result, s)
		}
	}
	return result
}

// ServerExists checks if a server with the given name exists
func ServerExists(servers []models.Server, name string) bool {
	return FindServerByName(servers, name) != nil
}

// SiteExists checks if a site with the given site ID exists on a server
func SiteExists(server *models.Server, siteID string) bool {
	return FindSiteBySiteID(server, siteID) != nil
}

// ParseSSLExpiry parses SSL certificate expiry date from openssl output format
// Input format: "Mar 15 12:00:00 2024 GMT" or similar
// Returns nil if parsing fails
func ParseSSLExpiry(expiryStr string) *time.Time {
	// Try common formats from openssl x509 -enddate output
	formats := []string{
		"Jan 2 15:04:05 2006 MST",
		"Jan  2 15:04:05 2006 MST",
		"Jan 02 15:04:05 2006 MST",
		"2 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, expiryStr); err == nil {
			return &t
		}
	}

	return nil
}
