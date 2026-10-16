# elastic_log_clean

Go rewrite of the Elastic diagnostic / cluster / server log cleaner, using the same omit + obfuscate + `report.yaml` model as [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean).

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

## Download

Pre-built binaries are on the [Releases](https://github.com/nwlterry/elastic_log_clean/releases) page (`v1.2.1`).

| Archive | Use on |
|---------|--------|
| `elastic_log_clean_1.2.1_linux_amd64.tar.gz` | RHEL / most servers |
| `elastic_log_clean_1.2.1_linux_arm64.tar.gz` | Linux aarch64 |
| `elastic_log_clean_1.2.1_darwin_amd64.tar.gz` | Intel macOS |
| `elastic_log_clean_1.2.1_darwin_arm64.tar.gz` | Apple Silicon |
| `elastic_log_clean_1.2.1_windows_amd64.zip` | Windows |

```bash
curl -fsSL -O https://github.com/nwlterry/elastic_log_clean/releases/download/v1.2.1/elastic_log_clean_1.2.1_linux_amd64.tar.gz
tar -xzf elastic_log_clean_1.2.1_linux_amd64.tar.gz
chmod +x elastic_log_clean
./elastic_log_clean --version
./elastic_log_clean -i diagnostics.zip -o scrubbed-diagnostics.zip
```

Each archive includes the binary plus `elastic_default.yaml` and `elastic_ip_name_map.yaml`.

## Windows 11

Unpack `elastic_log_clean_1.2.1_windows_amd64.zip`, then from PowerShell:

```powershell
cd ~\Downloads\elastic_log_clean_1.2.1_windows_amd64
powershell -ExecutionPolicy Bypass -File .\elastic_log_clean.ps1 -InputPath .\local-diagnostics.zip
.\elastic_log_clean.ps1 -InputPath C:\logs\elasticsearch -VerboseLog
.\elastic_log_clean.ps1 -InputPath .\cluster.log -Config .\elastic_ip_name_map.yaml
```

`-VerboseLog` prints every cleaned/omitted file. The exe can also be run directly:

```powershell
.\elastic_log_clean.exe -v -i .\local-diagnostics.zip -o .\scrubbed-local-diagnostics.zip
```

## Build

Requires Go 1.22+.

```bash
make          # writes bin/elastic_log_clean
make test
make install  # PREFIX=/usr/local
bash scripts/elastic_log_clean.sh --version
```

## Usage

```bash
./bin/elastic_log_clean -i local-diagnostics.zip -o scrubbed-local-diagnostics.zip
./bin/elastic_log_clean -i /var/log/elasticsearch -o /tmp/es-logs-cleaned
./bin/elastic_log_clean -c config/elastic_default.yaml -i prod.log -o prod.cleaned.log -r /tmp/report.yaml -w 8
./bin/elastic_log_clean -c config/elastic_ip_name_map.yaml -i cluster.zip -o scrubbed-cluster.zip
echo "node 10.20.30.41 password=secret" | ./bin/elastic_log_clean --stdin
./bin/elastic_log_clean --print-default-config
```

Flags match must-gather-clean: `-c` `-i` `-o` `-w` `-r`.

Do **not** share `report.yaml`.

## Matcher guards (v1.2.1)

Email and IPv6 candidates are filtered after the regex so diagnostic noise is not tokenized:

- Java / Maven coords such as `io.netty.transport@4.1.135.Final` stay as-is (TLD `Final` / version-shaped domain).
- Clock values such as `07:29:19` and `2026-09-24T07:29:19,123` are not treated as IPv6.
- Real addresses (`sre@corp.example.com`, `2001:db8::1`) are still tokenized.

## Config: IP only vs IP + custom names

Passing `-c` **replaces** the built-in config. It does not merge. A file that only contains `type: IP` will rewrite IPs and skip MAC / email / secrets / PEM / omit.

Default (`type: IP`, `replacementType: Consistent`) assigns tokens:

`10.99.1.8` → `x-ipv4-0000000001-x` (same source IP always gets the same token). Loopback `127.0.0.1` / `0.0.0.0` / `::1` is left alone.

The `IP` rule has no name map. Pin specific IPs or hostnames with `Keywords`. Keywords run **before** IP, so a mapped address is not tokenized again.

Copy [config/elastic_ip_name_map.yaml](config/elastic_ip_name_map.yaml) and edit the `replacement` maps:

```yaml
obfuscate:
  - type: IP
    replacementType: Consistent
    target: All
  - type: Keywords
    target: FileContents
    replacement:
      10.20.30.41: node-a
      10.20.30.42: node-b
      es-prod-01: node-a
      es-prod-02: node-b
  - type: Keywords
    target: FilePath
    replacement:
      es-prod-01: node-a
      es-prod-02: node-b
```

On a line `es-prod-01 10.20.30.41 10.99.1.8`:

| original      | result                |
|---------------|-----------------------|
| `es-prod-01`  | `node-a`              |
| `10.20.30.41` | `node-a`              |
| `10.99.1.8`   | `x-ipv4-0000000001-x` |

`FilePath` rewrites directory/file names; `FileContents` rewrites log text.

Do not commit a filled-in map or `report.yaml`.

## Layout

```
cmd/elastic_log_clean/main.go
pkg/clean/
config/elastic_default.yaml
config/elastic_ip_name_map.yaml
examples/sample-elasticsearch.log
scripts/elastic_log_clean.sh
Makefile
```

See [GROUP.md](GROUP.md).

## Principal / username redaction

In addition to IPs, emails, MACs, LDAP DNs and certificate subjects, the cleaner
redacts Kafka-style and MDS-style user identities while preserving the surrounding
structure:

| form | example | result |
|------|---------|--------|
| `User:` principal | `Principal = User:alice` | `Principal = User:x-user-…-x` |
| dotted name | `User:first_name.last_name` | `User:x-user-…-x` |
| `principal:User:` | `principal:User:bob` | `principal:User:x-user-…-x` |
| `username=` | `username=carol` / `username="carol"` | `username=x-user-…-x` |
| `userPrincipalName=` | `userPrincipalName=dave` | `userPrincipalName=x-user-…-x` |
| JSON principal | `"principal":"eve"` | `"principal":"x-user-…-x"` |
| login phrase | `login successful for user frank` | `login successful for user x-user-…-x` |
| MDS request log | `[date] INFO host - alice [utc] GET` | `host - x-user-…-x [utc] GET` |

The same consistent token is used for the same name across the run. Already-redacted
tokens (`x-…`) are left alone.
