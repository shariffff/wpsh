package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	wordmonDir = ".wordmon"
	ansibleDir = "ansible"
)

// GetWordmonDir returns the path to ~/.wordmon/
func GetWordmonDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, wordmonDir)
}

// GetAnsibleDir returns the path to ~/.wordmon/ansible/
func GetAnsibleDir() string {
	return filepath.Join(GetWordmonDir(), ansibleDir)
}

// IsInitialized checks if ~/.wordmon/ansible/ exists and has content
func IsInitialized() bool {
	ansiblePath := GetAnsibleDir()

	// Check if directory exists
	info, err := os.Stat(ansiblePath)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check if it has content (look for provision.yml)
	provisionFile := filepath.Join(ansiblePath, "provision.yml")
	_, err = os.Stat(provisionFile)
	return err == nil
}

// DetectAnsibleSource finds the ansible directory in the repository
// It looks for ansible/ relative to the CLI binary location or common paths
func DetectAnsibleSource() (string, error) {
	// Try relative path (development mode)
	relPaths := []string{
		"../ansible",    // from cli/ directory
		"ansible",       // from root directory
		"../../ansible", // from cli/cmd/ or deeper
	}

	for _, relPath := range relPaths {
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			continue
		}

		// Check if this path exists and has provision.yml
		provisionFile := filepath.Join(absPath, "provision.yml")
		if _, err := os.Stat(provisionFile); err == nil {
			return absPath, nil
		}
	}

	// Try system install path (for future package installations)
	systemPath := "/usr/local/share/wordmon/ansible"
	if _, err := os.Stat(filepath.Join(systemPath, "provision.yml")); err == nil {
		return systemPath, nil
	}

	return "", fmt.Errorf("could not find ansible source directory")
}

// Initialize sets up the ~/.wordmon directory and copies ansible files
func Initialize() error {
	wordmonPath := GetWordmonDir()
	ansiblePath := GetAnsibleDir()

	// Create ~/.wordmon/ directory
	if err := os.MkdirAll(wordmonPath, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", wordmonPath, err)
	}

	// Find ansible source
	ansibleSource, err := DetectAnsibleSource()
	if err != nil {
		return fmt.Errorf("failed to locate ansible directory: %w", err)
	}

	// Check if ansible directory already exists
	if _, err := os.Stat(ansiblePath); err == nil {
		return fmt.Errorf("ansible directory already exists at %s", ansiblePath)
	}

	// Copy ansible directory
	if err := copyDir(ansibleSource, ansiblePath); err != nil {
		return fmt.Errorf("failed to copy ansible files: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// GetAnsiblePath returns the path to use for ansible playbooks
// Checks in order: ~/.wordmon/ansible/, /usr/local/share/wordmon/ansible/, relative path
func GetAnsiblePath() (string, error) {
	// First check user's local copy
	userPath := GetAnsibleDir()
	if _, err := os.Stat(filepath.Join(userPath, "provision.yml")); err == nil {
		return userPath, nil
	}

	// Check system install
	systemPath := "/usr/local/share/wordmon/ansible"
	if _, err := os.Stat(filepath.Join(systemPath, "provision.yml")); err == nil {
		return systemPath, nil
	}

	// Fall back to detecting source (development mode)
	return DetectAnsibleSource()
}
