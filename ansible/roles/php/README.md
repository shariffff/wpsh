# PHP Role

Installs and configures PHP-FPM for WordPress hosting.

## What It Does

- Adds the ondrej/php PPA repository
- Installs PHP and common extensions
- Configures PHP-FPM with security hardening
- Installs Composer and WP-CLI

## Variables

All variables are defined in `defaults/main.yml` and can be overridden in `group_vars/all.yml` or via `--extra-vars`.

| Variable | Default | Description |
|----------|---------|-------------|
| `php_version` | `"8.3"` | PHP version to install |
| `php_extensions` | (list) | PHP extensions to install |
| `php_upload_max_filesize` | `"64M"` | Maximum upload file size |
| `php_post_max_size` | `"64M"` | Maximum POST data size |
| `php_memory_limit` | `"256M"` | Memory limit per script |
| `php_max_execution_time` | `300` | Max script execution time (seconds) |
| `php_max_input_vars` | `3000` | Maximum input variables |
| `php_fpm_pm` | `"ondemand"` | Process manager mode |
| `php_disabled_functions` | (list) | Functions disabled for security |
| `php_install_composer` | `true` | Install Composer globally |
| `php_install_wpcli` | `true` | Install WP-CLI globally |

## Security

The `php_disabled_functions` list disables functions commonly exploited in WordPress attacks:
- Shell execution: `exec`, `shell_exec`, `system`, `passthru`
- Process control: `pcntl_*`, `proc_*`, `popen`
- System info: `disk_free_space`, `posix_*`

Remove functions from this list only if your application requires them.

## Example Override

```yaml
# group_vars/all.yml
php_version: "8.2"
php_memory_limit: "512M"
php_upload_max_filesize: "128M"
```

## Handlers

- `restart php-fpm` - Restarts PHP-FPM service after configuration changes
