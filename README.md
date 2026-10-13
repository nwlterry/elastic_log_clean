# elastic_log_clean

Go log cleaner for Elasticsearch diagnostic bundles, cluster logs, and server logs. Same omit + obfuscate + `report.yaml` model as [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean).

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

## Download

Pre-built binaries are on the [Releases](https://github.com/nwlterry/elastic_log_clean/releases) page (`v1.0.0`).

| Archive | Use on |
|---------|--------|
| `elastic_log_clean_1.0.0_linux_amd64.tar.gz` | RHEL / most servers |
| `elastic_log_clean_1.0.0_linux_arm64.tar.gz` | Linux aarch64 |
| `elastic_log_clean_1.0.0_darwin_amd64.tar.gz` | Intel macOS |
| `elastic_log_clean_1.0.0_darwin_arm64.tar.gz` | Apple Silicon |
| `elastic_log_clean_1.0.0_windows_amd64.zip` | Windows 11 |

```bash
curl -fsSL -O https://github.com/nwlterry/elastic_log_clean/releases/download/v1.0.0/elastic_log_clean_1.0.0_linux_amd64.tar.gz
tar -xzf elastic_log_clean_1.0.0_linux_amd64.tar.gz
chmod +x elastic_log_clean
./elastic_log_clean --version
./elastic_log_clean -v -i diagnostics.zip -o scrubbed-diagnostics.zip
```

## Windows 11

```powershell
powershell -ExecutionPolicy Bypass -File .\elastic_log_clean.ps1 -InputPath .\diagnostics.zip
.\elastic_log_clean.ps1 -InputPath C:\logs\elasticsearch -VerboseLog
.\elastic_log_clean.ps1 -InputPath .\cluster.log -Config .\elastic_ip_name_map.yaml
```

## What is redacted

Consistent tokens (`x-ipv4-…`, `x-email-…`, `x-ldapdn-…`, `x-certsubj-…`, `x-redacted-secret-x`). Same source value always maps to the same token.

- IPv4 / IPv6 / MAC
- Email addresses
- LDAP DNs (`bind_dn`, `base_dn`) and search filters (`sAMAccountName=…`)
- Certificate subjects and issuers (`subject=`, `issuer=`, OpenSSL `/C=…/CN=…`)
- Passwords, API keys, Basic/Bearer auth, PEM blocks, URL userinfo
- Optional IP-to-name map via `config/elastic_ip_name_map.yaml`

## Guards

- `io.netty.transport@4.1.135.Final` is not an email
- `07:29:19` is not IPv6
- `x-redacted-secret-x@host` left by URL credential redaction is not reclassified as an email
- LDAP DNs keep multi-word values (`OU=Service Accounts`) and the full `DC=` chain, including a trailing `}}`

Do **not** share `report.yaml`.
