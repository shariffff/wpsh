# Database Role

Installs and configures MariaDB with security hardening.

## What It Does

- Installs MariaDB server and client
- Applies performance and security configuration
- Creates `wpshbot` admin user for site management
- Removes test database and anonymous users

## Variables

All variables are defined in `defaults/main.yml`.

| Variable                          | Default                | Description                    |
| --------------------------------- | ---------------------- | ------------------------------ |
| `mariadb_performance_schema`      | `false`                | Enable performance schema      |
| `mariadb_binary_logging`          | `false`                | Enable binary logging          |
| `mariadb_innodb_buffer_pool_size` | `"256M"`               | InnoDB buffer pool size        |
| `mariadb_max_connections`         | `100`                  | Maximum connections            |
| `mariadb_character_set`           | `"utf8mb4"`            | Default character set          |
| `mariadb_collation`               | `"utf8mb4_unicode_ci"` | Default collation              |
| `mariadb_slow_query_log`          | `false`                | Enable slow query logging      |
| `mariadb_long_query_time`         | `2`                    | Slow query threshold (seconds) |

## Required Variables

Set in `group_vars/all.yml`:

```yaml
mysql_wpshbot_password: 'your-secure-password'
```

## Tuning Guidelines

**`mariadb_innodb_buffer_pool_size`:**

- Small VPS (1-2GB RAM): `"128M"` to `"256M"`
- Medium VPS (4GB RAM): `"512M"` to `"1G"`
- Dedicated DB server: 50-70% of available RAM

**`mariadb_binary_logging`:**

- Enable (`true`) if you need replication or point-in-time recovery
- Disable (`false`) for single-server setups to save disk space

## Handlers

- `restart mariadb` - Restarts MariaDB after configuration changes
- `daemon-reload` - Reloads systemd after override changes
