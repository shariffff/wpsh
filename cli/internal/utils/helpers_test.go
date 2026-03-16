package utils

import (
	"testing"

	"github.com/wpsh/cli/pkg/models"
)

func TestFindServerByName(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", IP: "1.1.1.1"},
		{Name: "server2", IP: "2.2.2.2"},
		{Name: "server3", IP: "3.3.3.3"},
	}

	tests := []struct {
		name       string
		serverName string
		wantNil    bool
		wantIP     string
	}{
		{"find first server", "server1", false, "1.1.1.1"},
		{"find middle server", "server2", false, "2.2.2.2"},
		{"find last server", "server3", false, "3.3.3.3"},
		{"not found", "server4", true, ""},
		{"empty name", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindServerByName(servers, tt.serverName)
			if tt.wantNil && result != nil {
				t.Errorf("FindServerByName() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindServerByName() = nil, want server with IP %s", tt.wantIP)
			}
			if !tt.wantNil && result != nil && result.IP != tt.wantIP {
				t.Errorf("FindServerByName() IP = %v, want %v", result.IP, tt.wantIP)
			}
		})
	}
}

func TestFindServerIndexByName(t *testing.T) {
	servers := []models.Server{
		{Name: "server1"},
		{Name: "server2"},
		{Name: "server3"},
	}

	tests := []struct {
		name       string
		serverName string
		wantIndex  int
	}{
		{"find first", "server1", 0},
		{"find middle", "server2", 1},
		{"find last", "server3", 2},
		{"not found", "server4", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindServerIndexByName(servers, tt.serverName)
			if result != tt.wantIndex {
				t.Errorf("FindServerIndexByName() = %v, want %v", result, tt.wantIndex)
			}
		})
	}
}

func TestFindSiteBySiteID(t *testing.T) {
	server := &models.Server{
		Name: "testserver",
		Sites: []models.Site{
			{SiteID: "site1", PrimaryDomain: "site1.com"},
			{SiteID: "site2", PrimaryDomain: "site2.com"},
		},
	}

	tests := []struct {
		name       string
		server     *models.Server
		siteID     string
		wantNil    bool
		wantDomain string
	}{
		{"find first site", server, "site1", false, "site1.com"},
		{"find second site", server, "site2", false, "site2.com"},
		{"not found", server, "site3", true, ""},
		{"nil server", nil, "site1", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteBySiteID(tt.server, tt.siteID)
			if tt.wantNil && result != nil {
				t.Errorf("FindSiteBySiteID() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindSiteBySiteID() = nil, want site")
			}
			if !tt.wantNil && result != nil && result.PrimaryDomain != tt.wantDomain {
				t.Errorf("FindSiteBySiteID() domain = %v, want %v", result.PrimaryDomain, tt.wantDomain)
			}
		})
	}
}

func TestFindSiteIndexBySiteID(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{SiteID: "site1"},
			{SiteID: "site2"},
		},
	}

	tests := []struct {
		name      string
		server    *models.Server
		siteID    string
		wantIndex int
	}{
		{"find first", server, "site1", 0},
		{"find second", server, "site2", 1},
		{"not found", server, "site3", -1},
		{"nil server", nil, "site1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteIndexBySiteID(tt.server, tt.siteID)
			if result != tt.wantIndex {
				t.Errorf("FindSiteIndexBySiteID() = %v, want %v", result, tt.wantIndex)
			}
		})
	}
}

func TestFindSiteByDomain(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{
				SiteID:        "site1",
				PrimaryDomain: "primary.com",
				Domains: []models.Domain{
					{Domain: "primary.com"},
					{Domain: "www.primary.com"},
				},
			},
			{
				SiteID:        "site2",
				PrimaryDomain: "other.com",
				Domains: []models.Domain{
					{Domain: "other.com"},
				},
			},
		},
	}

	tests := []struct {
		name       string
		server     *models.Server
		domain     string
		wantNil    bool
		wantSiteID string
	}{
		{"find by primary domain", server, "primary.com", false, "site1"},
		{"find by additional domain", server, "www.primary.com", false, "site1"},
		{"find other site", server, "other.com", false, "site2"},
		{"not found", server, "notfound.com", true, ""},
		{"nil server", nil, "primary.com", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteByDomain(tt.server, tt.domain)
			if tt.wantNil && result != nil {
				t.Errorf("FindSiteByDomain() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindSiteByDomain() = nil, want site")
			}
			if !tt.wantNil && result != nil && result.SiteID != tt.wantSiteID {
				t.Errorf("FindSiteByDomain() siteID = %v, want %v", result.SiteID, tt.wantSiteID)
			}
		})
	}
}

func TestGetProvisionedServers(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", Status: "provisioned"},
		{Name: "server2", Status: "unprovisioned"},
		{Name: "server3", Status: "provisioned"},
		{Name: "server4", Status: "error"},
	}

	result := GetProvisionedServers(servers)

	if len(result) != 2 {
		t.Errorf("GetProvisionedServers() returned %d servers, want 2", len(result))
	}

	for _, s := range result {
		if s.Status != "provisioned" {
			t.Errorf("GetProvisionedServers() included server with status %s", s.Status)
		}
	}
}

func TestServerExists(t *testing.T) {
	servers := []models.Server{
		{Name: "server1"},
		{Name: "server2"},
	}

	tests := []struct {
		name       string
		serverName string
		want       bool
	}{
		{"exists", "server1", true},
		{"exists 2", "server2", true},
		{"not exists", "server3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ServerExists(servers, tt.serverName)
			if result != tt.want {
				t.Errorf("ServerExists() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestSiteExists(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{SiteID: "site1"},
			{SiteID: "site2"},
		},
	}

	tests := []struct {
		name   string
		server *models.Server
		siteID string
		want   bool
	}{
		{"exists", server, "site1", true},
		{"exists 2", server, "site2", true},
		{"not exists", server, "site3", false},
		{"nil server", nil, "site1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SiteExists(tt.server, tt.siteID)
			if result != tt.want {
				t.Errorf("SiteExists() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestParseSSLExpiry(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNil   bool
		wantYear  int
		wantMonth int
		wantDay   int
	}{
		{
			name:      "standard openssl format",
			input:     "Mar 15 12:00:00 2024 GMT",
			wantNil:   false,
			wantYear:  2024,
			wantMonth: 3,
			wantDay:   15,
		},
		{
			name:      "single digit day",
			input:     "Jan 5 08:30:00 2025 GMT",
			wantNil:   false,
			wantYear:  2025,
			wantMonth: 1,
			wantDay:   5,
		},
		{
			name:      "padded single digit day",
			input:     "Jan  5 08:30:00 2025 GMT",
			wantNil:   false,
			wantYear:  2025,
			wantMonth: 1,
			wantDay:   5,
		},
		{
			name:      "double digit day",
			input:     "Dec 25 23:59:59 2026 GMT",
			wantNil:   false,
			wantYear:  2026,
			wantMonth: 12,
			wantDay:   25,
		},
		{
			name:    "invalid format",
			input:   "2024-03-15",
			wantNil: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantNil: true,
		},
		{
			name:    "garbage input",
			input:   "not a date at all",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSSLExpiry(tt.input)
			if tt.wantNil {
				if result != nil {
					t.Errorf("ParseSSLExpiry(%q) = %v, want nil", tt.input, result)
				}
				return
			}
			if result == nil {
				t.Errorf("ParseSSLExpiry(%q) = nil, want non-nil", tt.input)
				return
			}
			if result.Year() != tt.wantYear {
				t.Errorf("ParseSSLExpiry(%q) year = %d, want %d", tt.input, result.Year(), tt.wantYear)
			}
			if int(result.Month()) != tt.wantMonth {
				t.Errorf("ParseSSLExpiry(%q) month = %d, want %d", tt.input, result.Month(), tt.wantMonth)
			}
			if result.Day() != tt.wantDay {
				t.Errorf("ParseSSLExpiry(%q) day = %d, want %d", tt.input, result.Day(), tt.wantDay)
			}
		})
	}
}
