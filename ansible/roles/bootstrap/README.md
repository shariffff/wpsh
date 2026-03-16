# Bootstrap Role

Prepares the server with base packages and the WPSH system user.

## What It Does

- Creates the `wpsh` system user for administration
- Installs required system packages
- Removes unnecessary packages (snapd, lxd)
- Configures base services (fail2ban, redis, certbot)

## Variables

All variables are defined in `defaults/main.yml`.

| Variable             | Default | Description         |
| -------------------- | ------- | ------------------- |
| `required_packages`  | (list)  | Packages to install |
| `packages_to_remove` | (list)  | Packages to remove  |

## Default Packages

**Installed:**

- `acl` - Access control lists for file permissions
- `gnupg` - GPG for package verification
- `cron` - Task scheduling
- `fail2ban` - Intrusion prevention
- `redis-server` - In-memory caching
- `certbot` - Let's Encrypt SSL
- `unattended-upgrades` - Automatic security updates

**Removed:**

- `lxd`, `lxcfs`, `snapd` - Container/snap packages (not needed)

## Required Variables

Set in `group_vars/all.yml`:

```yaml
wpsh_ssh_key: '~/.ssh/wpsh.pub' # or paste key directly
```

## Example: Add Custom Package

```yaml
# group_vars/all.yml
required_packages:
  - acl
  - gnupg
  - cron
  - fail2ban
  - redis-server
  - certbot
  - htop # Add monitoring tool
```
