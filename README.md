# WordMon

Automated WordPress hosting on Ubuntu servers. One command to provision, one command to deploy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shariffff/wordmon/main/install.sh | bash
```

Then initialize:

```bash
wordmon init
```

## Quick Start

```bash
# 1. Add your server
wordmon server add

# 2. Provision it (installs Nginx, PHP, MariaDB, Redis, SSL)
wordmon server provision myserver

# 3. Create a WordPress site
wordmon site create

# 4. Issue SSL certificate
wordmon domain ssl
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
wordmon server add          # Add a server
wordmon server provision    # Provision server with LEMP stack
wordmon server list         # List servers

wordmon site create         # Create WordPress site
wordmon site list           # List sites
wordmon site delete         # Delete site

wordmon domain add          # Add domain to site
wordmon domain ssl          # Issue SSL certificate

wordmon config show         # Show configuration
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
git clone https://github.com/shariffff/wordmon.git
cd wordmon
make build    # Build CLI
make test     # Run tests
```

## License

MIT
