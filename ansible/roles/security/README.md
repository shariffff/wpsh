# Security Role

Configures UFW firewall and SSH hardening for server security.

## What It Does

- Installs and configures UFW firewall
- Applies SSH hardening settings
- Sets default deny policy for incoming connections

## Variables

All variables are defined in `defaults/main.yml`.

| Variable | Default | Description |
|----------|---------|-------------|
| `ufw_allowed_ports` | `[22, 80, 443]` | Ports to allow through firewall |
| `ufw_default_policy` | `"deny"` | Default policy for incoming |
| `ssh_permit_root_login` | `"prohibit-password"` | Root login policy |
| `ssh_password_authentication` | `"no"` | Allow password auth |
| `ssh_pubkey_authentication` | `"yes"` | Allow pubkey auth |
| `ssh_max_auth_tries` | `3` | Max auth attempts |
| `ssh_login_grace_time` | `60` | Seconds to authenticate |

## SSH Options Explained

**`ssh_permit_root_login`:**
- `"no"` - Completely disable root login
- `"yes"` - Allow root login (not recommended)
- `"prohibit-password"` - Root can only login with SSH keys

## Example: Add Custom Port

```yaml
# group_vars/all.yml
ufw_allowed_ports:
  - 22
  - 80
  - 443
  - 8080  # Custom application port
```

## Handlers

- `restart ssh` - Restarts SSH service after configuration changes
