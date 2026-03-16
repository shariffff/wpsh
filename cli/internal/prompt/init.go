package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// InitInput holds the input for init setup
type InitInput struct {
	SSHPublicKey string
	CertbotEmail string
}

// PromptInitSetup prompts for initial setup configuration
func PromptInitSetup() (*InitInput, error) {
	input := &InitInput{}

	fmt.Println()
	fmt.Println("Let's configure some one-time settings for your WPSH installation.")
	fmt.Println()

	// SSH public key selection
	sshPubKeys, err := findSSHPublicKeys()
	if err != nil || len(sshPubKeys) == 0 {
		// Fallback to manual input if no keys found
		homeDir, _ := os.UserHomeDir()
		defaultKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa.pub")

		keyPrompt := &survey.Input{
			Message: "SSH public key file (for wp-sh user):",
			Default: defaultKeyPath,
			Help:    "Path to the SSH public key that will be authorized for the wp-sh user on servers",
		}
		if err := survey.AskOne(keyPrompt, &input.SSHPublicKey, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
	} else if len(sshPubKeys) == 1 {
		// Auto-select when only one key is available
		input.SSHPublicKey = sshPubKeys[0]
		fmt.Printf("Using SSH public key: %s\n", sshPubKeys[0])
	} else {
		// Show picker with available keys
		options := append(sshPubKeys, "Enter path manually")
		keyPrompt := &survey.Select{
			Message: "SSH public key (for wp-sh user):",
			Options: options,
			Help:    "Select the SSH public key to authorize for the wp-sh user on servers",
		}
		var selectedKey string
		if err := survey.AskOne(keyPrompt, &selectedKey); err != nil {
			return nil, err
		}

		if selectedKey == "Enter path manually" {
			homeDir, _ := os.UserHomeDir()
			defaultKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa.pub")
			keyPrompt := &survey.Input{
				Message: "SSH public key file path:",
				Default: defaultKeyPath,
			}
			if err := survey.AskOne(keyPrompt, &input.SSHPublicKey, survey.WithValidator(survey.Required)); err != nil {
				return nil, err
			}
		} else {
			input.SSHPublicKey = selectedKey
		}
	}

	// Certbot email
	emailPrompt := &survey.Input{
		Message: "Email for Let's Encrypt SSL certificates:",
		Help:    "This email will receive certificate expiration notices from Let's Encrypt",
	}
	if err := survey.AskOne(emailPrompt, &input.CertbotEmail, survey.WithValidator(survey.Required), survey.WithValidator(validateEmail)); err != nil {
		return nil, err
	}

	return input, nil
}

// findSSHPublicKeys looks for public SSH keys in ~/.ssh/
func findSSHPublicKeys() ([]string, error) {
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

		// Only include .pub files
		if !strings.HasSuffix(name, ".pub") {
			continue
		}

		keyPath := filepath.Join(sshDir, name)
		keys = append(keys, keyPath)
	}

	return keys, nil
}

// validateEmail validates an email address
func validateEmail(val interface{}) error {
	email, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid input")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("please enter a valid email address")
	}

	return nil
}
