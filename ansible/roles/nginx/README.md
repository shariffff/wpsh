# Nginx Role

Installs and configures Nginx from the official repository.

## What It Does

- Adds official Nginx repository
- Installs Nginx with optimal configuration
- Generates default self-signed SSL certificate
- Configures global settings for WordPress hosting

## Variables

All variables are defined in `defaults/main.yml`.

| Variable | Default | Description |
|----------|---------|-------------|
| `nginx_ppa` | (dynamic) | Nginx official repository URL |
| `nginx_worker_connections` | `8000` | Max connections per worker |

## Configuration Details

The role configures Nginx with:
- Official Nginx mainline packages (not Ubuntu's older version)
- Optimized worker settings
- Default SSL certificate for immediate HTTPS support
- Security headers and best practices

## Tuning Guidelines

**`nginx_worker_connections`:**
- Default `8000` is suitable for most VPS
- High-traffic sites: increase to `16000` or higher
- Formula: `worker_processes * worker_connections = max concurrent connections`

## Example Override

```yaml
# group_vars/all.yml
nginx_worker_connections: 16000
```

## Handlers

- `restart nginx` - Restarts Nginx after configuration changes

## Generated Files

- `/etc/nginx/nginx.conf` - Main configuration
- `/etc/nginx/ssl/default.*` - Default self-signed SSL certificate
