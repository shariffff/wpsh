package models

import "time"

// Domain represents a domain associated with a site
type Domain struct {
	Domain        string     `yaml:"domain" validate:"required,fqdn"`
	SSLEnabled    bool       `yaml:"ssl_enabled"`
	SSLIssuedAt   *time.Time `yaml:"ssl_issued_at,omitempty"`
	SSLExpiresAt  *time.Time `yaml:"ssl_expires_at,omitempty"`
}
