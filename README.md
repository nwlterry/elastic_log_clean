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
