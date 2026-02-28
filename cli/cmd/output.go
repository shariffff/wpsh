package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// CommandResult represents a JSON response for command execution
type CommandResult struct {
	Success bool                   `json:"success"`
	Action  string                 `json:"action,omitempty"`
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// isJSONOutput checks if the command should output JSON
func isJSONOutput(cmd *cobra.Command) bool {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	return jsonFlag
}

// outputSuccess outputs a success message, either as JSON or human-readable
func outputSuccess(cmd *cobra.Command, action string, data map[string]interface{}) {
	if isJSONOutput(cmd) {
		result := CommandResult{
			Success: true,
			Action:  action,
			Data:    data,
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	} else {
		// Human-readable output
		switch action {
		case "server_added":
			color.Green("✓ Server '%s' added successfully", data["name"])
		case "server_provisioned":
			color.Green("✓ Server '%s' provisioned successfully", data["name"])
		case "server_removed":
			color.Green("✓ Server '%s' removed from inventory", data["name"])
		case "server_updated":
			color.Green("✓ Server '%s' updated successfully", data["name"])
		case "server_healthy":
			color.Green("✓ Server '%s' is healthy", data["name"])
		case "site_created":
			color.Green("✓ WordPress site created successfully")
		case "site_deleted":
			color.Green("✓ Site '%s' deleted successfully", data["domain"])
		case "domain_added":
			color.Green("✓ Domain '%s' added successfully", data["domain"])
		case "domain_removed":
			color.Green("✓ Domain '%s' removed successfully", data["domain"])
		case "ssl_issued":
			color.Green("✓ SSL certificate issued successfully")
		default:
			color.Green("✓ Operation completed successfully")
		}
	}
}

// outputError outputs an error message, either as JSON or human-readable
func outputError(cmd *cobra.Command, message string, err error) {
	if isJSONOutput(cmd) {
		result := CommandResult{
			Success: false,
			Message: message,
			Error:   err.Error(),
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	} else {
		color.Red("Error: %s: %v", message, err)
	}
}

// outputInfo outputs an informational message (only in non-JSON mode)
func outputInfo(cmd *cobra.Command, format string, args ...interface{}) {
	if !isJSONOutput(cmd) {
		fmt.Printf(format, args...)
	}
}
