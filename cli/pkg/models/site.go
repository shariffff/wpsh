package models

import (
	"time"

	"gopkg.in/yaml.v3"
)

// Database holds database connection info for a site
type Database struct {
	Name string `yaml:"name" validate:"required"`
	User string `yaml:"user" validate:"required"`
	Host string `yaml:"host" validate:"required"`
}

// Metadata holds additional site information
type Metadata struct {
	BackupEnabled bool       `yaml:"backup_enabled"`
	LastBackup    *time.Time `yaml:"last_backup,omitempty"`
}

// Site represents a WordPress site on a server
type Site struct {
	SiteID        string    `yaml:"site_id" validate:"required,alphanum"`
	PrimaryDomain string    `yaml:"primary_domain" validate:"required,fqdn"`
	CreatedAt     time.Time `yaml:"created_at"`
	AdminUser     string    `yaml:"admin_user" validate:"required"`
	AdminEmail    string    `yaml:"admin_email" validate:"required,email"`
	Domains       []Domain  `yaml:"domains"`
	Database      Database  `yaml:"database"`
	PHPVersion    string    `yaml:"php_version"`
	Metadata      Metadata  `yaml:"metadata"`
	Notes         string    `yaml:"notes,omitempty"`
}

// rawSite is used for YAML unmarshalling with backwards compatibility
type rawSite struct {
	SiteID        string    `yaml:"site_id"`
	SystemName    string    `yaml:"system_name"` // Legacy field for backwards compatibility
	PrimaryDomain string    `yaml:"primary_domain"`
	CreatedAt     time.Time `yaml:"created_at"`
	AdminUser     string    `yaml:"admin_user"`
	AdminEmail    string    `yaml:"admin_email"`
	Domains       []Domain  `yaml:"domains"`
	Database      Database  `yaml:"database"`
	PHPVersion    string    `yaml:"php_version"`
	Metadata      Metadata  `yaml:"metadata"`
	Notes         string    `yaml:"notes,omitempty"`
}

// UnmarshalYAML implements custom unmarshalling for backwards compatibility
// It reads both "site_id" (new) and "system_name" (legacy) fields
func (s *Site) UnmarshalYAML(value *yaml.Node) error {
	var raw rawSite
	if err := value.Decode(&raw); err != nil {
		return err
	}

	// Prefer site_id, fall back to system_name for backwards compatibility
	if raw.SiteID != "" {
		s.SiteID = raw.SiteID
	} else if raw.SystemName != "" {
		s.SiteID = raw.SystemName
	}

	s.PrimaryDomain = raw.PrimaryDomain
	s.CreatedAt = raw.CreatedAt
	s.AdminUser = raw.AdminUser
	s.AdminEmail = raw.AdminEmail
	s.Domains = raw.Domains
	s.Database = raw.Database
	s.PHPVersion = raw.PHPVersion
	s.Metadata = raw.Metadata
	s.Notes = raw.Notes

	return nil
}
