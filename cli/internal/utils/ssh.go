package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"github.com/wp-sh/cli/pkg/models"
)

// TestSSHConnection tests SSH connectivity to a server
func TestSSHConnection(server models.Server) error {
	// Expand home directory in key file path
	keyFile := server.SSH.KeyFile
	if strings.HasPrefix(keyFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		keyFile = filepath.Join(homeDir, keyFile[1:])
	}

	// Read SSH private key
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read SSH key file %s: %w", keyFile, err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	// Configure SSH client with TOFU host key verification
	// This validates against known_hosts if the file exists and the host is known,
	// or automatically accepts and saves unknown host keys
	config := &ssh.ClientConfig{
		User: server.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: trustOnFirstUseCallback(),
		Timeout:         10 * time.Second,
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", server.IP, server.SSH.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed to %s: %w", addr, err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Test command execution
	output, err := session.CombinedOutput("echo 'wp-sh-test'")
	if err != nil {
		return fmt.Errorf("test command failed: %w", err)
	}

	if strings.TrimSpace(string(output)) != "wp-sh-test" {
		return fmt.Errorf("unexpected test output: %s", output)
	}

	return nil
}

// getHostKeyCallback returns a host key callback using the user's known_hosts file
func getHostKeyCallback() (ssh.HostKeyCallback, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	return knownhosts.New(knownHostsPath)
}

// trustOnFirstUseCallback returns a callback that accepts any host key
// and adds it to known_hosts on first connection (TOFU model)
func trustOnFirstUseCallback() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// If we can't get home dir, just accept the key
			return nil
		}

		knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

		// Try to read existing known_hosts
		callback, err := knownhosts.New(knownHostsPath)
		if err == nil {
			// File exists, check against it
			err = callback(hostname, remote, key)
			if err == nil {
				return nil // Key is known and matches
			}
			// Check if it's a KeyError (unknown host or key mismatch)
			if keyErr, ok := err.(*knownhosts.KeyError); ok {
				// If Want is not empty, it means we expected different keys (mismatch)
				if len(keyErr.Want) > 0 {
					return fmt.Errorf("host key mismatch for %s - possible security issue", hostname)
				}
				// Want is empty, so host is unknown - fall through to add it
			}
		}

		// Key not in known_hosts, add it (TOFU)
		// Ensure .ssh directory exists
		sshDir := filepath.Join(homeDir, ".ssh")
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return nil // Accept key even if we can't save it
		}

		// Append to known_hosts
		f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil // Accept key even if we can't save it
		}
		defer f.Close()

		// Format the known_hosts line
		line := knownhosts.Line([]string{hostname}, key)
		if _, err := f.WriteString(line + "\n"); err != nil {
			return nil // Accept key even if we can't save it
		}

		return nil
	}
}
