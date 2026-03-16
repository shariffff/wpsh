# WPSH Ansible Playbooks

Ansible playbooks and roles for automated WordPress hosting infrastructure.

## Overview

This directory contains all Ansible automation for:

- Server provisioning (LEMP stack)
- WordPress site deployment
- Domain and SSL management
- Site operations (deletion, management)

## Quick Start

### Prerequisites

```bash
# Install Ansible
pip install ansible

# Install required collections
ansible-galaxy install -r requirements.yml
```

### Core Playbooks

**Provision a fresh Ubuntu server:**

```bash
ansible-playbook provision.yml -i "SERVER_IP," -u root
```

**Create a WordPress site:**

```bash
ansible-playbook website.yml -i "SERVER_IP," -u wpsh \
  --extra-vars "domain=example.com system_name=examplecom wp_admin_user=admin wp_admin_email=admin@example.com wp_admin_password=SecurePass123"
```

**Domain management:**

```bash
# Add domain
ansible-playbook playbooks/domain_management.yml -i "IP," -u wpsh \
  --extra-vars "operation=add_domain domain=newdomain.com system_name=sitename"

# Remove domain
ansible-playbook playbooks/domain_management.yml -i "IP," -u wpsh \
  --extra-vars "operation=remove_domain domain=olddomain.com"

# Issue SSL certificate
ansible-playbook playbooks/domain_management.yml -i "IP," -u wpsh \
  --extra-vars "operation=issue_ssl domain=example.com certbot_email=admin@example.com"
```

**Delete a site:**

```bash
ansible-playbook playbooks/delete_site.yml -i "IP," -u wpsh \
  --extra-vars "system_name=examplecom"
```

## Playbook Overview

| Playbook                          | Purpose                 | User   | Required Variables               |
| --------------------------------- | ----------------------- | ------ | -------------------------------- |
| `provision.yml`                   | Full server setup       | `root` | See group_vars/all.yml           |
| `website.yml`                     | Create WordPress site   | `wpsh` | domain, system*name, wp_admin*\* |
| `playbooks/domain_management.yml` | Add/remove domains, SSL | `wpsh` | operation, domain                |
| `playbooks/delete_site.yml`       | Remove site completely  | `wpsh` | system_name                      |

## Roles Architecture

| Role           | Purpose              | Key Tasks                                                                              |
| -------------- | -------------------- | -------------------------------------------------------------------------------------- |
| **bootstrap**  | Base system setup    | Creates wpsh user, installs base packages, certbot, fail2ban, redis                    |
| **database**   | MariaDB installation | Installs MariaDB, creates wpshbot admin user, secures installation                     |
| **nginx**      | Web server setup     | Installs Nginx from official repo, configures global settings, generates default SSL   |
| **php**        | PHP installation     | Installs PHP 8.3 from ondrej/php PPA, configures PHP-FPM, installs Composer and WP-CLI |
| **security**   | Security hardening   | Configures UFW firewall (ports 22/80/443), SSH hardening                               |
| **website**    | Site deployment      | Creates site user, database, PHP-FPM pool, Nginx vhost, installs WordPress             |
| **libs**       | Reusable tasks       | add_domain, remove_domain, issue_ssl, cleanup_server                                   |
| **operations** | Server operations    | delete_site, manage_database, manage_domain, manage_systemd, verify_connection         |

## Required Variables

### Global Variables (group_vars/all.yml)

| Variable                 | Description                                             | Required |
| ------------------------ | ------------------------------------------------------- | -------- |
| `wpsh_ssh_key`           | SSH public key for wpsh user (file path or key content) | Yes      |
| `mysql_wpshbot_password` | MySQL admin password                                    | Yes      |
| `certbot_email`          | Email for Let's Encrypt                                 | Yes      |

### Website Creation (website.yml)

| Variable            | Description              | Example             |
| ------------------- | ------------------------ | ------------------- |
| `domain`            | Primary domain name      | `example.com`       |
| `system_name`       | System identifier        | `examplecom`        |
| `wp_admin_user`     | WordPress admin username | `admin`             |
| `wp_admin_email`    | WordPress admin email    | `admin@example.com` |
| `wp_admin_password` | WordPress admin password | `SecurePass123`     |

## Server Directory Structure

After provisioning and site creation:

```
/sites/example.com/
├── public/            # WordPress root (web accessible)
└── logs/              # Site-specific logs

/etc/nginx/sites-available/example.com/
├── example.com        # Main server configuration
├── server/            # Server block includes
├── location/          # Location block includes
├── before/            # Pre-processing rules
└── after/             # Post-processing rules (redirects)

/etc/php/8.3/fpm/pool.d/
└── examplecom.conf    # Dedicated PHP-FPM pool

/home/examplecom/      # Site user home directory
```

## Running with Tags

Execute specific parts of playbooks:

```bash
# Run only bootstrap tasks
ansible-playbook provision.yml -i "SERVER_IP," -u root --tags bootstrap

# Available tags for provision.yml
--tags bootstrap    # Base system setup
--tags database     # MariaDB installation
--tags nginx        # Nginx setup
--tags php          # PHP installation
--tags security     # Security hardening

# Available tags for domain_management.yml
--tags add_domain    # Add domain only
--tags remove_domain # Remove domain only
--tags issue_ssl     # Issue SSL certificate only
```

## Configuration Files

### ansible.cfg

```ini
[defaults]
roles_path = ./roles
```

### group_vars/all.yml

Global variables applied to all hosts. Set required variables here or pass via `--extra-vars`.

### inventory/

Example inventory files for different environments.

## Execution Flow

### provision.yml

```
bootstrap → database → nginx → php → security
```

### website.yml

```
Pre-tasks (generate credentials) →
website role:
  users → php → nginx → files → database → wordpress → cron
```

### domain_management.yml

```
Route based on operation variable →
  add_domain | remove_domain | issue_ssl →
  Reload Nginx
```

## Best Practices

### Inventory Format

```bash
# Ad-hoc single host (note the comma)
ansible-playbook playbook.yml -i "192.168.1.100," -u user

# Inventory file
ansible-playbook playbook.yml -i inventory/production -u user
```

### Variable Precedence

1. Command line `--extra-vars` (highest)
2. Playbook `vars:` section
3. `group_vars/all.yml`
4. Role defaults (lowest)

### Idempotency

All playbooks are designed to be idempotent - running them multiple times produces the same result without side effects.

### Security Notes

- Never commit sensitive variables to git
- Use Ansible Vault for secrets: `ansible-vault encrypt_string 'secret' --name 'variable_name'`
- Store passwords in environment variables or vault files

## Technology Stack

- **Target OS**: Ubuntu 24.04 LTS
- **Web Server**: Nginx (official repository)
- **PHP**: 8.3 from ondrej/php PPA
- **Database**: MariaDB
- **Cache**: Redis
- **SSL**: Let's Encrypt via Certbot
- **Security**: UFW, Fail2ban

## Troubleshooting

### Check Syntax

```bash
ansible-playbook --syntax-check provision.yml
```

### Dry Run

```bash
ansible-playbook provision.yml -i "IP," -u root --check
```

### Verbose Output

```bash
ansible-playbook provision.yml -i "IP," -u root -vvv
```

### Test Connection

```bash
ansible all -i "IP," -u wpsh -m ping
```

## Advanced Usage

### Custom PHP Version

```bash
# Modify group_vars/all.yml or override
ansible-playbook website.yml -i "IP," -u wpsh \
  --extra-vars "domain=example.com system_name=examplecom ... php_version=8.2"
```

### Skip Tags

```bash
# Skip security hardening
ansible-playbook provision.yml -i "IP," -u root --skip-tags security
```

### Limit Hosts

```bash
# When using inventory files
ansible-playbook provision.yml -i inventory/production --limit webserver1
```

## Development

### Testing Changes

```bash
# Syntax check
ansible-playbook --syntax-check provision.yml

# Dry run
ansible-playbook provision.yml -i "IP," -u root --check

# Run from repository root
cd .. && make test-ansible
```

### Adding New Roles

1. Create role structure: `mkdir -p roles/newrole/{tasks,handlers,templates,defaults,files}`
2. Add tasks in `roles/newrole/tasks/main.yml`
3. Add handlers in `roles/newrole/handlers/main.yml`
4. Include role in appropriate playbook

## Support

For issues specific to Ansible playbooks:

1. Check playbook syntax: `ansible-playbook --syntax-check playbook.yml`
2. Run with verbose output: `-vvv`
3. Review logs on target server: `/var/log/syslog`, site-specific logs in `/sites/*/logs/`

For general WPSH help, see the [main README](../README.md) or use the CLI: `wpsh --help`
