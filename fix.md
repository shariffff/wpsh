|-------------------------------|-----------------------------------------------|------------------------|
| Unrestricted NOPASSWD sudo | roles/bootstrap/tasks/wordmon-user.yml:32 | Full system compromise |
| Plaintext DB credentials | roles/database/templates/_.j2 | Credential exposure |
| _.\*:ALL,GRANT DB privileges | roles/database/tasks/main.yml:57 | Database compromise |
| Excessive ignore_errors: true | roles/website/tasks/wordpress.yml:20,34,41,48 | Silent failures |
