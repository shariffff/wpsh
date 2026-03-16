package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wp-sh/cli/internal/utils"
	"github.com/wp-sh/cli/pkg/models"
)

// ServerInput holds the input for server creation
type ServerInput struct {
	Name     string
	Hostname string
	IP       string
	SSHUser  string
	SSHPort  int
	SSHKey   string
}

// PromptServerAdd prompts for server details
func PromptServerAdd() (*ServerInput, error) {
	input := &ServerInput{}

	// Server name
	namePrompt := &survey.Input{
		Message: "Server name:",
		Help:    "A friendly name to identify this server (e.g., production-1, staging)",
	}
	if err := survey.AskOne(namePrompt, &input.Name, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// IP address
	ipPrompt := &survey.Input{
		Message: "IP address:",
		Help:    "The server's IP address from your cloud provider",
	}
	if err := survey.AskOne(ipPrompt, &input.IP, survey.WithValidator(survey.Required), survey.WithValidator(utils.ValidateIP)); err != nil {
		return nil, err
	}

	// Set hostname to IP by default
	input.Hostname = input.IP

	// SSH key selection
	sshKeys, err := findSSHKeys()
	if err != nil || len(sshKeys) == 0 {
		// Fallback to manual input if no keys found
		homeDir, _ := os.UserHomeDir()
		defaultKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa")

		keyPrompt := &survey.Input{
			Message: "SSH private key file:",
			Default: defaultKeyPath,
			Help:    "Path to the SSH private key for authentication",
		}
		if err := survey.AskOne(keyPrompt, &input.SSHKey, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
	} else if len(sshKeys) == 1 {
		// Auto-select when only one key is available
		input.SSHKey = sshKeys[0]
		fmt.Printf("Using SSH key: %s\n", sshKeys[0])
	} else {
		// Show picker with available keys
		keyPrompt := &survey.Select{
			Message: "SSH private key:",
			Options: sshKeys,
			Help:    "Select the SSH key to use for authentication",
		}
		var selectedKey string
		if err := survey.AskOne(keyPrompt, &selectedKey); err != nil {
			return nil, err
		}
		input.SSHKey = selectedKey
	}

	// SSH user
	userPrompt := &survey.Input{
		Message: "SSH user:",
		Default: "root",
		Help:    "The SSH user to connect as (use 'root' for initial provisioning)",
	}
	if err := survey.AskOne(userPrompt, &input.SSHUser, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// SSH port
	portPrompt := &survey.Input{
		Message: "SSH port:",
		Default: "22",
		Help:    "The SSH port number",
	}
	var portStr string
	if err := survey.AskOne(portPrompt, &portStr, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}
	fmt.Sscanf(portStr, "%d", &input.SSHPort)
	if input.SSHPort == 0 {
		input.SSHPort = 22
	}

	// Confirmation
	if err := confirmServerAdd(input); err != nil {
		return nil, err
	}

	return input, nil
}

// findSSHKeys looks for private SSH keys in ~/.ssh/
func findSSHKeys() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip public keys, known_hosts, config, etc.
		if strings.HasSuffix(name, ".pub") ||
			name == "known_hosts" ||
			name == "known_hosts.old" ||
			name == "config" ||
			name == "authorized_keys" ||
			name == "environment" {
			continue
		}

		// Check if it looks like a private key by reading first line
		keyPath := filepath.Join(sshDir, name)
		content, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}

		// Private keys typically start with "-----BEGIN"
		if strings.HasPrefix(string(content), "-----BEGIN") {
			keys = append(keys, keyPath)
		}
	}

	return keys, nil
}

// ToServer converts ServerInput to models.Server
func (si *ServerInput) ToServer() models.Server {
	return models.Server{
		Name:     si.Name,
		Hostname: si.Hostname,
		IP:       si.IP,
		SSH: models.SSHConfig{
			User:    si.SSHUser,
			Port:    si.SSHPort,
			KeyFile: si.SSHKey,
		},
		Status: "unprovisioned",
		Sites:  []models.Site{},
	}
}

func confirmServerAdd(input *ServerInput) error {
	fmt.Println("\nServer Configuration:")
	fmt.Printf("  Name:     %s\n", input.Name)
	fmt.Printf("  IP:       %s\n", input.IP)
	fmt.Printf("  SSH Key:  %s\n", input.SSHKey)
	fmt.Printf("  SSH User: %s\n", input.SSHUser)
	fmt.Printf("  SSH Port: %d\n", input.SSHPort)

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Proceed with provisioning?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("cancelled")
	}

	return nil
}
