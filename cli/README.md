# WordMon CLI

A command-line interface tool for managing WordPress hosting infrastructure using Ansible. WordMon provides both an intuitive, interactive mode for manual operations and a script mode with flags for automation and CI/CD pipelines.

## Features

- **Dual Operating Modes**: Interactive prompts for manual use, script mode with flags for automation
- **YAML-based State Management**: All configuration stored in `~/.wordmon/wordmon.yaml`
- **Ansible Integration**: Seamlessly executes existing Ansible playbooks
- **Server Management**: Add, list, remove, and provision servers
- **Site Management**: Create, list, and delete WordPress sites
- **Domain Management**: Add domains and manage SSL certificates

## Installation

### Prerequisites

- Go 1.24 or higher (or Docker for building without Go)
- Ansible installed and configured
- SSH access to target servers

### Build from Source

```bash
# Clone or navigate to the repository
cd /path/to/ansible/cli

# Build the binary (requires Go 1.24+)
make build

# Or build using Docker (no Go required)
make docker-build

# Install to /usr/local/bin (requires sudo)
make install

# Or install to ~/bin (no sudo required)
make install-user
```

### Verify Installation

```bash
wordmon version
wordmon --help
```

## Quick Start

### 1. Initialize Configuration

```bash
wordmon init
```

This creates `~/.wordmon/wordmon.yaml` with default settings.

### 2. Configure Ansible Path

Edit `~/.wordmon/wordmon.yaml` and set the correct Ansible project path:

```yaml
ansible:
  path: '/Users/yourname/Projects/ansible' # Update this path
  roles_path: './roles'
  inventory_path: '/tmp/wordmon-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'
```

### 3. Add a Server

```bash
wordmon server add
```

Follow the interactive prompts to add server details:

- Server name (e.g., production-1)
- Hostname or IP address
- SSH user and port
- SSH private key file

### 4. List Servers

```bash
wordmon server list
```

### 5. Validate Configuration

```bash
wordmon config validate
```

## Operating Modes

WordMon CLI supports two modes of operation:

### Interactive Mode (Default)

When you run commands without flags, the CLI guides you through the process with interactive prompts.

```bash
wordmon site create
# Prompts you for: server, domain, admin credentials
# (site ID is auto-generated from domain)
```

**Use interactive mode when:**

- Learning the tool
- Performing manual operations
- You want validation and helpful hints
- Exploring available options

### Script Mode (Non-Interactive)

Provide all parameters as command-line flags for fully automated operations.

```bash
wordmon site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!
# --site-id is optional (auto-generated from domain if not provided)
```

**Use script mode when:**

- Automating deployments
- Running in CI/CD pipelines
- Scripting repetitive tasks
- No user interaction is possible

**Common flags for script mode:**

- `--non-interactive`: Required flag to enable script mode
- `--force`: Skip confirmation prompts
- `--skip-ssh-check`: Skip SSH connectivity validation

## Commands

### Configuration Management

```bash
# Show current configuration
wordmon config show

# Validate configuration
wordmon config validate

# Edit configuration in your preferred editor
wordmon config edit
```

### Server Management

```bash
# Add a new server
wordmon server add

# List all servers
wordmon server list

# Remove a server
wordmon server remove <name>

# Provision a server
wordmon server provision <name>

# Provision with options
wordmon server provision <name> --force              # Skip confirmation
wordmon server provision <name> --skip-ssh-check     # Skip SSH connectivity test
```

### Site Management

```bash
# Create a new WordPress site (interactive)
wordmon site create

# Create a site non-interactively (site-id auto-generated)
wordmon site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# Create with explicit site-id
wordmon site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --site-id mysite \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# List all sites
wordmon site list

# List sites on a specific server
wordmon site list --server production-1

# Delete a site (interactive selection)
wordmon site delete

# Delete a specific site (by site ID)
wordmon site delete --server production-1 --site mysiteid

# Force delete without confirmation
wordmon site delete --server production-1 --site mysiteid --force
```

### Domain Management

```bash
# Add a domain to a site (interactive)
wordmon domain add

# Add domain with automatic SSL
# (prompts will ask if you want to issue SSL)

# Remove a domain (interactive selection)
wordmon domain remove

# Force remove without confirmation
wordmon domain remove --force

# Issue SSL certificate for a domain (interactive)
wordmon domain ssl

# The CLI will:
# - Show only domains without SSL
# - Prompt for Let's Encrypt email
# - Obtain and configure SSL certificate
# - Update Nginx to use HTTPS
# - Track SSL expiration in configuration
```

## Configuration File

The configuration file is located at `~/.wordmon/wordmon.yaml`. Here's an example structure:

```yaml
version: '1.0'

ansible:
  path: '/Users/sharif/Projects/ansible'
  roles_path: './roles'
  inventory_path: '/tmp/wordmon-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'

global_vars:
  certbot_email: 'admin@example.com'
  mysql_wordmonbot_password: '${MYSQL_WORDMONBOT_PASSWORD}'
  wordmon_ssh_key: '~/.ssh/wordmon_rsa.pub'

servers:
  - name: 'production-1'
    hostname: 'prod1.example.com'
    ip: '203.0.113.10'
    ssh:
      user: 'wordmon'
      port: 22
      key_file: '~/.ssh/wordmon_rsa'
    status: 'unprovisioned'
    sites: []
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Format Code

```bash
make fmt
```

### Clean Build Artifacts

```bash
make clean
```

## Project Structure

```
cli/
├── cmd/                  # Command definitions
│   ├── root.go          # Root command
│   ├── config.go        # Config commands
│   ├── server.go        # Server commands
│   └── version.go       # Version command
├── internal/
│   ├── config/          # Configuration management
│   ├── ansible/         # Ansible integration (coming soon)
│   ├── state/           # State management (coming soon)
│   ├── prompt/          # Interactive prompts
│   └── utils/           # Utilities (coming soon)
├── pkg/
│   └── models/          # Data models
├── templates/           # Templates (inventory, etc.)
├── main.go             # Entry point
├── Makefile            # Build automation
└── README.md           # This file
```

## Roadmap

- [ ] Shell completion scripts
- [ ] Comprehensive error handling
- [ ] Installation script
- [ ] Release automation
