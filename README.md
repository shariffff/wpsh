# WPSH

Automated WordPress hosting on Ubuntu servers. One command to provision, one command to deploy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shariffff/wp-sh/main/install.sh | bash
```

Then initialize:

```bash
wp-sh init
```

## Quick Start

```bash
# 1. Add your server
wp-sh server add

# 2. Provision it (installs Nginx, PHP, MariaDB, Redis, SSL)
wp-sh server provision myserver

# 3. Create a WordPress site
wp-sh site create

# 4. Issue SSL certificate
wp-sh domain ssl
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
wp-sh server add          # Add a server
wp-sh server provision    # Provision server with LEMP stack
wp-sh server list         # List servers

wp-sh site create         # Create WordPress site
wp-sh site list           # List sites
wp-sh site delete         # Delete site

wp-sh domain add          # Add domain to site
wp-sh domain ssl          # Issue SSL certificate

wp-sh config show         # Show configuration
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
git clone https://github.com/shariffff/wp-sh.git
cd wp-sh
make build    # Build CLI
make test     # Run tests
```

## License

MIT
