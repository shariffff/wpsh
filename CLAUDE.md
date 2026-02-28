# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WordMon is an automated WordPress hosting infrastructure using Ansible. It consists of:

- **Ansible playbooks** (`ansible/`) for server provisioning and WordPress site management
- **Go CLI tool** (`cli/`) that provides an interactive wrapper around Ansible playbooks

## Repository Structure

```
wordmon/
├── ansible/          # All Ansible playbooks and roles
├── cli/              # Go CLI tool
├── docs/             # Additional documentation
├── Makefile          # Top-level build automation
├── README.md         # Project overview
├── CLAUDE.md         # This file
└── version.txt       # Version tracking
```

## Running Ansible Playbooks Directly

**Note**: Users typically interact via the CLI (`wordmon` command), but you can run Ansible directly for testing or advanced usage.

### Prerequisites

```bash
# Install Ansible dependencies
cd ansible
ansible-galaxy install -r requirements.yml
```

### Core Playbooks

**Provision a fresh Ubuntu server:**

```bash
cd ansible
ansible-playbook provision.yml -i "SERVER_IP," -u root
```

**Create a WordPress site:**

```bash
cd ansible
ansible-playbook website.yml -i "SERVER_IP," -u wordmon \
  --extra-vars "domain=example.com site_id=examplecom wp_admin_user=admin wp_admin_email=admin@example.com wp_admin_password=SecurePass123"
```

**Domain management:**

```bash
cd ansible
# Add domain
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordmon \
  --extra-vars "operation=add_domain domain=newdomain.com site_id=sitename"

# Remove domain
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordmon \
  --extra-vars "operation=remove_domain domain=olddomain.com"

# Issue SSL
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordmon \
  --extra-vars "operation=issue_ssl domain=example.com certbot_email=admin@example.com"
```

**Delete a site:**

```bash
cd ansible
ansible-playbook playbooks/delete_site.yml -i "IP," -u wordmon \
  --extra-vars "site_id=examplecom"
```

### Running with Specific Tags

```bash
cd ansible
# Run only bootstrap tasks
ansible-playbook provision.yml -i "SERVER_IP," -u root --tags bootstrap

# Available tags: bootstrap, database, nginx, php, security, website
```

## Building and Testing

### Top-level Commands (Recommended)

```bash
# From repository root

# Build the CLI
make build

# Run all tests (CLI + Ansible syntax check)
make test

# Format Go code
make fmt

# Lint Go code
make lint

# Install to /usr/local/bin (requires sudo)
make install

# Install to ~/bin (no sudo)
make install-user

# Clean build artifacts
make clean

# Test Ansible playbooks
make test-ansible
```

### CLI-specific Commands

```bash
cd cli

# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Install to /usr/local/bin
make install

# Install to ~/bin (no sudo)
make install-user

# Clean build artifacts
make clean
```

### CLI Commands Overview

```bash
# Initialization (run once after install)
wordmon init

# Configuration management
wordmon config show
wordmon config validate
wordmon config edit

# Server management
wordmon server add
wordmon server list
wordmon server remove <name>
wordmon server provision <name>

# Site management
wordmon site create
wordmon site list
wordmon site delete

# Domain management
wordmon domain add
wordmon domain remove
wordmon domain ssl
```

## Architecture

### Ansible Structure (ansible/ directory)

**Playbook Execution Flow:**

- `ansible/provision.yml` → Runs 5 roles sequentially: bootstrap → database → nginx → php → security
- `ansible/website.yml` → Runs website role, which imports 7 task files in order: users → php → nginx → files → database → wordpress → cron
- `ansible/playbooks/domain_management.yml` → Routes to libs role tasks based on `operation` variable

**Roles Architecture:**

- **bootstrap**: Creates `wordmon` system user, installs base packages, certbot, fail2ban, redis
- **database**: Installs MariaDB, creates `wordmonbot` admin user
- **nginx**: Installs from official repo, sets up global config, generates default SSL
- **php**: Installs PHP 8.3 from ondrej/php PPA, configures PHP-FPM pools, installs Composer and WP-CLI
- **security**: Configures UFW (ports 22/80/443), SSH hardening
- **website**: Orchestrates site creation through 7 task files (users, php, nginx, files, database, wordpress, cron)
- **libs**: Reusable task libraries (add_domain, remove_domain, issue_ssl, cleanup_server)
- **operations**: Server operation tasks (delete_site, manage_database, manage_domain, manage_systemd, verify_connection)

**Variable Flow:**

- Global variables defined in `ansible/group_vars/all.yml`
- Playbook-specific variables set in playbook `vars:` or `pre_tasks:`
- Runtime variables passed via `--extra-vars`
- `ansible/website.yml` generates random credentials in pre_tasks and sets dynamic facts

**Server Directory Layout:**

```
/sites/example.com/
├── public/            # WordPress root
└── logs/              # Site-specific logs

/etc/nginx/sites-available/example.com/
├── example.com        # Main server config
├── server/            # Server block includes
├── location/          # Location block includes
├── before/            # Pre-processing rules
└── after/             # Post-processing (redirects)
```

### CLI Structure (Go)

**Module Path:** `github.com/wordmon/cli`

**Package Organization:**

- `cmd/`: Cobra command definitions (root, init, config, server, site, domain, version)
- `internal/config/`: YAML config loading/saving, validation
- `internal/ansible/`: Ansible inventory generation and playbook execution
- `internal/installer/`: Setup logic for copying ansible files to ~/.wordmon/
- `internal/prompt/`: Interactive prompts using survey library
- `internal/state/`: State updates for servers/sites/domains
- `internal/utils/`: Validation utilities
- `pkg/models/`: Data models (Server, Site, Domain)
- `templates/`: Ansible inventory templates

**State Management:**

- All state stored in `~/.wordmon/wordmon.yaml`
- Config structure: version, ansible settings, global_vars, servers array
- Each server has: name, hostname, ip, ssh config, status, sites array
- Each site has: domain, site_id, admin credentials, domains array
- Each domain has: name, ssl_enabled, ssl_expiry

**Ansible Integration:**

- `wordmon init` copies ansible/ directory to `~/.wordmon/ansible/` on first run
- CLI detects ansible location (user's ~/.wordmon/ansible/ or system install path)
- CLI generates temporary inventory files at runtime
- Executes Ansible playbooks with real-time output streaming
- Updates state after successful operations
- SSH connectivity checks before provisioning

## Required Variables

Set in `ansible/group_vars/all.yml` or pass via `--extra-vars`:

| Variable                    | Description                                                |
| --------------------------- | ---------------------------------------------------------- |
| `wordmon_ssh_key`           | SSH public key for wordmon user (file path or key content) |
| `mysql_wordmonbot_password` | MySQL admin password                                       |
| `certbot_email`             | Email for Let's Encrypt                                    |

For `website.yml`:
| Variable | Description |
|----------|-------------|
| `domain` | Primary domain name |
| `site_id` | Site identifier (used for user, db, PHP pool) - auto-generated if not provided |
| `wp_admin_user` | WordPress admin username |
| `wp_admin_email` | WordPress admin email |
| `wp_admin_password` | WordPress admin password |

## Technology Stack

- **Target OS**: Ubuntu 24.04
- **Web Server**: Nginx (official repo)
- **PHP**: 8.3 (ondrej/php PPA)
- **Database**: MariaDB
- **Cache**: Redis
- **SSL**: Let's Encrypt (Certbot)
- **Security**: UFW, Fail2ban
- **CLI**: Go 1.21+, Cobra framework, Survey prompts

## Key Patterns

**Ansible Conventions:**

- Use `ansible.builtin.*` modules explicitly (not implicit `module_name`)
- Tasks in libs role are designed to be included via `ansible.builtin.include_role` with `tasks_from:`
- Handlers defined in `roles/*/handlers/main.yml` use `listen:` directive for deduplication
- All playbooks use `become: true` except where explicitly noted
- Inventory format: Use comma notation for ad-hoc hosts: `"SERVER_IP,"`

**CLI Conventions:**

- Commands support both interactive (prompts) and non-interactive (flags) modes
- State updates happen after successful Ansible execution
- Validation occurs before execution (config, SSH connectivity)
- Ansible paths resolved in order: `~/.wordmon/ansible/` → `/usr/local/share/wordmon/ansible/` → relative path (dev mode)
- `wordmon init` must be run once after installation to copy ansible files locally
