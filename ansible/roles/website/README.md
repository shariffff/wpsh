# Website Role

Creates and configures WordPress sites with isolated resources.

## What It Does

- Creates system user and group for the site
- Configures PHP-FPM pool with per-site isolation
- Sets up Nginx server configuration
- Creates site directories and sets permissions
- Creates MySQL database and user
- Installs and configures WordPress
- Sets up system cron for WP-Cron

## Variables

All variables are defined in `defaults/main.yml`.

| Variable | Default | Description |
|----------|---------|-------------|
| `php_version` | `"8.3"` | PHP version for the site |
| `site_php_pm` | `"dynamic"` | PHP-FPM process manager mode |
| `site_php_pm_max_children` | `5` | Max PHP-FPM workers |
| `site_php_pm_start_servers` | `1` | Initial workers (dynamic mode) |
| `site_php_pm_min_spare_servers` | `1` | Min idle workers |
| `site_php_pm_max_spare_servers` | `1` | Max idle workers |
| `site_php_pm_max_requests` | `500` | Requests before worker recycle |
| `wp_disable_cron` | `true` | Use system cron instead of WP-Cron |
| `wp_cron_interval` | `"*/5"` | Cron schedule (every 5 minutes) |
| `wp_debug` | `false` | Enable WordPress debug mode |

## Required Variables

Pass via `--extra-vars` when creating a site:

```yaml
domain: "example.com"
system_name: "examplecom"
wp_admin_user: "admin"
wp_admin_email: "admin@example.com"
wp_admin_password: "SecurePassword123"
```

## Process Manager Modes

**`dynamic`** (recommended for most sites):
- Adjusts workers based on load
- Good balance of performance and memory

**`ondemand`** (memory efficient):
- No workers at startup
- Spawns on demand, ideal for low-traffic sites

**`static`** (high performance):
- Fixed number of workers
- Best for high-traffic sites with predictable load

## Site Isolation

Each site runs with:
- Dedicated system user and group
- Isolated PHP-FPM pool
- Separate MySQL database and user
- Own directory structure under `/sites/{domain}/`

## Handlers

- `Reload nginx` - Reloads Nginx configuration
- `Reload php-fpm` - Reloads PHP-FPM configuration
