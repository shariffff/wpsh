package ansible

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/wpsh/cli/pkg/models"
)

// ExecutionResult holds the parsed results from Ansible output
type ExecutionResult struct {
	Ok      int
	Changed int
	Failed  int
}

// PlaybookResult holds the complete result from playbook execution
type PlaybookResult struct {
	Success   bool
	Output    []string
	DNSStatus *DNSStatus
	SSLInfo   *SSLInfo
}

// DNSStatus holds DNS check results parsed from Ansible output
type DNSStatus struct {
	Domain     string
	ResolvedIP string
	ServerIP   string
	Matches    bool
}

// SSLInfo holds SSL issuance results parsed from Ansible output
type SSLInfo struct {
	Domain string
	Expiry string
}

// Executor handles Ansible playbook execution
type Executor struct {
	ansiblePath  string
	invGenerator *InventoryGenerator
	verbose      bool
	dryRun       bool
	spinner      *spinner.Spinner
}

// NewExecutor creates a new Ansible executor
func NewExecutor(ansiblePath string) *Executor {
	return &Executor{
		ansiblePath:  ansiblePath,
		invGenerator: NewInventoryGenerator(),
		verbose:      false,
		dryRun:       false,
	}
}

// SetVerbose enables or disables verbose output
func (e *Executor) SetVerbose(verbose bool) {
	e.verbose = verbose
}

// SetDryRun enables or disables dry-run mode (--check in Ansible)
func (e *Executor) SetDryRun(dryRun bool) {
	e.dryRun = dryRun
}

// ExecutePlaybook runs an ansible-playbook command with the given parameters
func (e *Executor) ExecutePlaybook(playbookName string, server models.Server, extraVars map[string]interface{}, globalVars map[string]interface{}) error {
	// Expand home directory in ansible path if needed
	ansiblePath := e.ansiblePath
	if strings.HasPrefix(ansiblePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		ansiblePath = filepath.Join(homeDir, ansiblePath[1:])
	}

	// Generate inventory
	inventoryPath, err := e.invGenerator.Generate(server, fmt.Sprintf("wpsh %s", playbookName), globalVars)
	if err != nil {
		return fmt.Errorf("failed to generate inventory: %w", err)
	}
	defer e.invGenerator.Cleanup(inventoryPath)

	// Build playbook path
	playbookPath := filepath.Join(ansiblePath, playbookName)

	// Check if playbook exists
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		return fmt.Errorf("playbook not found: %s", playbookPath)
	}

	// Build command arguments
	args := []string{
		playbookPath,
		"-i", inventoryPath,
	}

	// Add verbose flag if enabled (only for ansible, not our spinner mode)
	if e.verbose {
		args = append(args, "-vv")
	}

	// Add dry-run flag if enabled
	if e.dryRun {
		args = append(args, "--check")
	}

	// Merge globalVars and extraVars for --extra-vars (highest precedence)
	// This ensures CLI-provided values override group_vars/all.yml
	allVars := make(map[string]interface{})
	for k, v := range globalVars {
		allVars[k] = v
	}
	for k, v := range extraVars {
		allVars[k] = v
	}

	// Add extra vars if any exist
	if len(allVars) > 0 {
		varsJSON, err := json.Marshal(allVars)
		if err != nil {
			return fmt.Errorf("failed to marshal extra vars: %w", err)
		}
		args = append(args, "--extra-vars", string(varsJSON))
	}

	// Create command
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = ansiblePath

	// Set environment variables
	cmd.Env = os.Environ()

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Use spinner mode (quiet) by default, verbose mode shows full output
	if e.verbose {
		// Verbose mode: show full Ansible output
		fmt.Printf("\n")
		color.Cyan("Running: ansible-playbook %s", strings.Join(args, " "))
		fmt.Printf("\n")

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start ansible-playbook: %w", err)
		}

		done := make(chan bool)
		go func() {
			e.streamOutput(stdout, false)
			done <- true
		}()
		go func() {
			e.streamOutput(stderr, true)
		}()
		<-done

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("ansible-playbook failed: %w", err)
		}
		return nil
	}

	// Spinner mode (default): show spinner with current task
	return e.executeWithSpinner(cmd, stdout, stderr)
}

// executeWithSpinner runs the command with a spinner showing current task
func (e *Executor) executeWithSpinner(cmd *exec.Cmd, stdout, stderr io.ReadCloser) error {
	// Initialize spinner
	e.spinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	e.spinner.Suffix = " Starting..."
	e.spinner.Start()

	if err := cmd.Start(); err != nil {
		e.spinner.Stop()
		return fmt.Errorf("failed to start ansible-playbook: %w", err)
	}

	// Buffers to store output
	var outputBuffer []string
	var errorBuffer []string
	var result ExecutionResult
	var currentTask string
	var failed bool
	var mu sync.Mutex

	// Regex patterns
	taskPattern := regexp.MustCompile(`^TASK \[(.+?)\]`)
	playPattern := regexp.MustCompile(`^PLAY \[(.+?)\]`)
	recapPattern := regexp.MustCompile(`ok=(\d+)\s+changed=(\d+).*failed=(\d+)`)
	failedPattern := regexp.MustCompile(`(FAILED!|fatal:)`)

	done := make(chan bool, 2)

	// Process stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			outputBuffer = append(outputBuffer, line)

			// Check for task name
			if matches := taskPattern.FindStringSubmatch(line); len(matches) > 1 {
				currentTask = matches[1]
				e.spinner.Suffix = " " + currentTask
			} else if matches := playPattern.FindStringSubmatch(line); len(matches) > 1 {
				e.spinner.Suffix = " " + matches[1]
			}

			// Check for failures
			if failedPattern.MatchString(line) {
				failed = true
			}

			// Parse recap
			if matches := recapPattern.FindStringSubmatch(line); len(matches) > 3 {
				fmt.Sscanf(matches[1], "%d", &result.Ok)
				fmt.Sscanf(matches[2], "%d", &result.Changed)
				fmt.Sscanf(matches[3], "%d", &result.Failed)
			}
			mu.Unlock()
		}
		done <- true
	}()

	// Process stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			errorBuffer = append(errorBuffer, line)
			if failedPattern.MatchString(line) {
				failed = true
			}
			mu.Unlock()
		}
		done <- true
	}()

	// Wait for both streams
	<-done
	<-done

	// Wait for command to finish
	cmdErr := cmd.Wait()
	e.spinner.Stop()

	// Show results
	if cmdErr != nil || failed || result.Failed > 0 {
		// Show failure
		color.Red("✗ Task failed: %s\n", currentTask)
		fmt.Println()

		// Show relevant error output (last 20 lines or lines containing errors)
		mu.Lock()
		e.printErrorContext(outputBuffer, errorBuffer)
		mu.Unlock()

		fmt.Println()
		color.Red("Failed: %d ok, %d changed, %d failed", result.Ok, result.Changed, result.Failed)
		if cmdErr != nil {
			return fmt.Errorf("ansible-playbook failed")
		}
		return fmt.Errorf("playbook completed with failures")
	}

	// Show success
	color.Green("✓ Completed: %d ok, %d changed, %d failed", result.Ok, result.Changed, result.Failed)
	return nil
}

// printErrorContext prints relevant lines from the output when an error occurs
func (e *Executor) printErrorContext(outputBuffer, errorBuffer []string) {
	// Print stderr if any
	for _, line := range errorBuffer {
		color.Red(line)
	}

	// Find and print lines around the failure
	failedPattern := regexp.MustCompile(`(?i)(FAILED|fatal:|TASK \[)`)
	inErrorContext := false
	contextLines := 0
	maxContextLines := 15

	for _, line := range outputBuffer {
		if failedPattern.MatchString(line) {
			inErrorContext = true
			contextLines = 0
		}

		if inErrorContext {
			if strings.Contains(line, "FAILED") || strings.Contains(line, "fatal:") {
				color.Red(line)
			} else if strings.Contains(line, "TASK [") {
				color.Cyan(line)
			} else {
				fmt.Println(line)
			}
			contextLines++
			if contextLines > maxContextLines && !strings.Contains(line, "fatal:") && !strings.Contains(line, "FAILED") {
				inErrorContext = false
			}
		}
	}
}

// ExecutePlaybookWithResult runs a playbook and returns parsed results
func (e *Executor) ExecutePlaybookWithResult(playbookName string, server models.Server, extraVars map[string]interface{}, globalVars map[string]interface{}) (*PlaybookResult, error) {
	// Expand home directory in ansible path if needed
	ansiblePath := e.ansiblePath
	if strings.HasPrefix(ansiblePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to expand home directory: %w", err)
		}
		ansiblePath = filepath.Join(homeDir, ansiblePath[1:])
	}

	// Generate inventory
	inventoryPath, err := e.invGenerator.Generate(server, fmt.Sprintf("wpsh %s", playbookName), globalVars)
	if err != nil {
		return nil, fmt.Errorf("failed to generate inventory: %w", err)
	}
	defer e.invGenerator.Cleanup(inventoryPath)

	// Build playbook path
	playbookPath := filepath.Join(ansiblePath, playbookName)

	// Check if playbook exists
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("playbook not found: %s", playbookPath)
	}

	// Build command arguments
	args := []string{
		playbookPath,
		"-i", inventoryPath,
	}

	// Add verbose flag if enabled
	if e.verbose {
		args = append(args, "-vv")
	}

	// Add dry-run flag if enabled
	if e.dryRun {
		args = append(args, "--check")
	}

	// Merge globalVars and extraVars for --extra-vars (highest precedence)
	// This ensures CLI-provided values override group_vars/all.yml
	allVars := make(map[string]interface{})
	for k, v := range globalVars {
		allVars[k] = v
	}
	for k, v := range extraVars {
		allVars[k] = v
	}

	// Add extra vars if any exist
	if len(allVars) > 0 {
		varsJSON, err := json.Marshal(allVars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal extra vars: %w", err)
		}
		args = append(args, "--extra-vars", string(varsJSON))
	}

	// Create command
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = ansiblePath
	cmd.Env = os.Environ()

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Execute with result capture
	return e.executeWithSpinnerAndResult(cmd, stdout, stderr)
}

// executeWithSpinnerAndResult runs the command with spinner and returns parsed results
func (e *Executor) executeWithSpinnerAndResult(cmd *exec.Cmd, stdout, stderr io.ReadCloser) (*PlaybookResult, error) {
	// Initialize spinner
	e.spinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	e.spinner.Suffix = " Starting..."
	e.spinner.Start()

	if err := cmd.Start(); err != nil {
		e.spinner.Stop()
		return nil, fmt.Errorf("failed to start ansible-playbook: %w", err)
	}

	// Buffers to store output
	var outputBuffer []string
	var errorBuffer []string
	var result ExecutionResult
	var currentTask string
	var failed bool
	var mu sync.Mutex

	// Regex patterns
	taskPattern := regexp.MustCompile(`^TASK \[(.+?)\]`)
	playPattern := regexp.MustCompile(`^PLAY \[(.+?)\]`)
	recapPattern := regexp.MustCompile(`ok=(\d+)\s+changed=(\d+).*failed=(\d+)`)
	failedPattern := regexp.MustCompile(`(FAILED!|fatal:)`)

	done := make(chan bool, 2)

	// Process stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			outputBuffer = append(outputBuffer, line)

			if matches := taskPattern.FindStringSubmatch(line); len(matches) > 1 {
				currentTask = matches[1]
				e.spinner.Suffix = " " + currentTask
			} else if matches := playPattern.FindStringSubmatch(line); len(matches) > 1 {
				e.spinner.Suffix = " " + matches[1]
			}

			if failedPattern.MatchString(line) {
				failed = true
			}

			if matches := recapPattern.FindStringSubmatch(line); len(matches) > 3 {
				fmt.Sscanf(matches[1], "%d", &result.Ok)
				fmt.Sscanf(matches[2], "%d", &result.Changed)
				fmt.Sscanf(matches[3], "%d", &result.Failed)
			}
			mu.Unlock()
		}
		done <- true
	}()

	// Process stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			errorBuffer = append(errorBuffer, line)
			if failedPattern.MatchString(line) {
				failed = true
			}
			mu.Unlock()
		}
		done <- true
	}()

	<-done
	<-done

	cmdErr := cmd.Wait()
	e.spinner.Stop()

	// Parse results
	playbookResult := &PlaybookResult{
		Success: cmdErr == nil && !failed && result.Failed == 0,
		Output:  outputBuffer,
	}

	// Parse DNS status and SSL info from output
	playbookResult.DNSStatus = parseDNSStatus(outputBuffer)
	playbookResult.SSLInfo = parseSSLInfo(outputBuffer)

	// Show results
	if cmdErr != nil || failed || result.Failed > 0 {
		color.Red("✗ Task failed: %s\n", currentTask)
		fmt.Println()
		mu.Lock()
		e.printErrorContext(outputBuffer, errorBuffer)
		mu.Unlock()
		fmt.Println()
		color.Red("Failed: %d ok, %d changed, %d failed", result.Ok, result.Changed, result.Failed)
		if cmdErr != nil {
			return playbookResult, fmt.Errorf("ansible-playbook failed")
		}
		return playbookResult, fmt.Errorf("playbook completed with failures")
	}

	color.Green("✓ Completed: %d ok, %d changed, %d failed", result.Ok, result.Changed, result.Failed)
	return playbookResult, nil
}

// parseDNSStatus parses DNS_STATUS line from Ansible output
func parseDNSStatus(output []string) *DNSStatus {
	// Pattern: DNS_STATUS: domain=example.com resolved_ip=1.2.3.4 server_ip=5.6.7.8 matches=true
	dnsPattern := regexp.MustCompile(`DNS_STATUS:\s*domain=(\S+)\s+resolved_ip=(\S+)\s+server_ip=(\S+)\s+matches=(\S+)`)

	for _, line := range output {
		if matches := dnsPattern.FindStringSubmatch(line); len(matches) > 4 {
			return &DNSStatus{
				Domain:     matches[1],
				ResolvedIP: matches[2],
				ServerIP:   matches[3],
				Matches:    matches[4] == "True" || matches[4] == "true",
			}
		}
	}
	return nil
}

// parseSSLInfo parses SSL_ISSUED line from Ansible output
func parseSSLInfo(output []string) *SSLInfo {
	// Pattern: SSL_ISSUED: domain=example.com expiry=Mar 15 12:00:00 2024 GMT
	sslPattern := regexp.MustCompile(`SSL_ISSUED:\s*domain=(\S+)\s+expiry=(.+)$`)

	for _, line := range output {
		if matches := sslPattern.FindStringSubmatch(line); len(matches) > 2 {
			return &SSLInfo{
				Domain: matches[1],
				Expiry: strings.TrimSpace(matches[2]),
			}
		}
	}
	return nil
}

// streamOutput reads and prints output with color coding
func (e *Executor) streamOutput(reader io.Reader, isError bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Color code based on content
		if isError {
			color.Red(line)
		} else if strings.Contains(line, "FAILED") || strings.Contains(line, "fatal:") {
			color.Red(line)
		} else if strings.Contains(line, "ok:") || strings.Contains(line, "skipping:") {
			color.Green(line)
		} else if strings.Contains(line, "changed:") {
			color.Yellow(line)
		} else if strings.Contains(line, "PLAY [") || strings.Contains(line, "TASK [") {
			color.Cyan(line)
		} else if strings.Contains(line, "PLAY RECAP") {
			color.Magenta(line)
		} else {
			fmt.Println(line)
		}
	}
}
