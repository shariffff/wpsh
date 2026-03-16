# WPSH

Automated WordPress hosting on Ubuntu servers. One command to provision, one command to deploy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shariffff/wpsh/main/install.sh | bash
```

Then initialize:

```bash
wpsh init
```

## Quick Start

```bash
# 1. Add your server
wpsh server add

# 2. Provision it (installs Nginx, PHP, MariaDB, Redis, SSL)
wpsh server provision myserver

# 3. Create a WordPress site
wpsh site create

# 4. Issue SSL certificate
wpsh domain ssl
```

## What It Does

**Server provisioning:**

- Nginx from official repo
- PHP 8.3 with optimized FPM pools
- MariaDB with secure defaults
- Redis for object caching
- Let's Encrypt SSL via Certbot
- UFW firewall + Fail2ban

**Site isolation:**

- Each site runs as its own Linux user
- Dedicated PHP-FPM pool per site
- Isolated file permissions

## Commands

```bash
wpsh server add          # Add a server
wpsh server provision    # Provision server with LEMP stack
wpsh server list         # List servers

wpsh site create         # Create WordPress site
wpsh site list           # List sites
wpsh site delete         # Delete site

wpsh domain add          # Add domain to site
wpsh domain ssl          # Issue SSL certificate

wpsh config show         # Show configuration
```

All commands support `--help` for details.

## Requirements

- Ansible 2.14+ on your local machine
- Ubuntu 24.04 target server with root SSH access

## Documentation

- [CLI Reference](cli/README.md)
- [Ansible Playbooks](ansible/README.md)

## Development

```bash
git clone https://github.com/shariffff/wpsh.git
cd wpsh
make build    # Build CLI
make test     # Run tests
```

## License

MIT
