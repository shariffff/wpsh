package ansible

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/wp-sh/cli/pkg/models"
)

//go:embed inventory.tmpl
var inventoryTemplate string

// InventoryData holds the data for inventory template
type InventoryData struct {
	Timestamp         string
	Server            models.Server
	Command           string
	PythonInterpreter string
	GlobalVars        map[string]string
}

// InventoryGenerator generates Ansible inventory files
type InventoryGenerator struct {
	outputDir string
}

// NewInventoryGenerator creates a new inventory generator
func NewInventoryGenerator() *InventoryGenerator {
	return &InventoryGenerator{
		outputDir: "/tmp",
	}
}

// Generate creates an inventory file for the given server
func (ig *InventoryGenerator) Generate(server models.Server, command string, globalVars map[string]interface{}) (string, error) {
	// Convert globalVars to string map
	varsMap := make(map[string]string)
	for key, val := range globalVars {
		varsMap[key] = fmt.Sprintf("%v", val)
	}

	// Expand environment variables in values
	for key, val := range varsMap {
		varsMap[key] = os.ExpandEnv(val)
	}

	// Get home directory once for reuse
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "" // Will skip home expansion if we can't get it
	}

	// Expand home directory in global vars (especially wp-sh_ssh_key)
	for key, val := range varsMap {
		if strings.HasPrefix(val, "~") && homeDir != "" {
			varsMap[key] = filepath.Join(homeDir, val[1:])
		}
	}

	// Expand home directory in SSH key file
	sshKeyFile := server.SSH.KeyFile
	if strings.HasPrefix(sshKeyFile, "~") && homeDir != "" {
		sshKeyFile = filepath.Join(homeDir, sshKeyFile[1:])
	}
	server.SSH.KeyFile = sshKeyFile

	// Prepare template data
	data := InventoryData{
		Timestamp:         time.Now().Format(time.RFC3339),
		Server:            server,
		Command:           command,
		PythonInterpreter: "/usr/bin/python3",
		GlobalVars:        varsMap,
	}

	// Parse template
	tmpl, err := template.New("inventory").Parse(inventoryTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse inventory template: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102-150405")
	outputPath := filepath.Join(ig.outputDir, fmt.Sprintf("wp-sh-%s-%s.ini", server.Name, timestamp))

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create inventory file: %w", err)
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("failed to execute inventory template: %w", err)
	}

	return outputPath, nil
}

// Cleanup removes a generated inventory file
func (ig *InventoryGenerator) Cleanup(inventoryPath string) error {
	if inventoryPath == "" {
		return nil
	}
	return os.Remove(inventoryPath)
}
