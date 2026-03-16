package models

import "time"

// SSHConfig holds SSH connection details for a server
type SSHConfig struct {
	User    string `yaml:"user" validate:"required"`
	Port    int    `yaml:"port" validate:"required,min=1,max=65535"`
	KeyFile string `yaml:"key_file" validate:"required"`
}

// ServerCredentials holds server-specific credentials
type ServerCredentials struct {
	MySQLWordmonbotPassword string `yaml:"mysql_wp-shbot_password,omitempty"`
}

// Server represents a managed server
type Server struct {
	Name          string            `yaml:"name" validate:"required"`
	Hostname      string            `yaml:"hostname" validate:"required"`
	IP            string            `yaml:"ip" validate:"required,ip"`
	SSH           SSHConfig         `yaml:"ssh"`
	Credentials   ServerCredentials `yaml:"credentials,omitempty"`
	Status        string            `yaml:"status" validate:"oneof=provisioned unprovisioned error"`
	ProvisionedAt *time.Time        `yaml:"provisioned_at,omitempty"`
	Sites         []Site            `yaml:"sites,omitempty"`
}
