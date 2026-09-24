# elastic_log_clean

Go rewrite of the Elastic diagnostic / cluster / server log cleaner, using the same omit + obfuscate + `report.yaml` model as [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean).

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

## Download

Pre-built binaries are on the [Releases](https://github.com/nwlterry/elastic_log_clean/releases) page (`v1.2.0`).

| Archive | Use on |
|---------|--------|
| `elastic_log_clean_1.2.0_linux_amd64.tar.gz` | RHEL / most servers |
| `elastic_log_clean_1.2.0_linux_arm64.tar.gz` | Linux aarch64 |
| `elastic_log_clean_1.2.0_darwin_amd64.tar.gz` | Intel macOS |
| `elastic_log_clean_1.2.0_darwin_arm64.tar.gz` | Apple Silicon |
| `elastic_log_clean_1.2.0_windows_amd64.zip` | Windows 11 |

```bash
curl -fsSL -O https://github.com/nwlterry/elastic_log_clean/releases/download/v1.2.0/elastic_log_clean_1.2.0_linux_amd64.tar.gz
tar -xzf elastic_log_clean_1.2.0_linux_amd64.tar.gz
chmod +x elastic_log_clean
./elastic_log_clean --version
./elastic_log_clean -v -i diagnostics.zip -o scrubbed-diagnostics.zip
```

Each archive includes the binary, YAML configs, and `elastic_log_clean.ps1`.

## Windows 11

Copy the Elastic Support Diagnostics zip or server log off the host, unpack the windows package, then:

```powershell
powershell -ExecutionPolicy Bypass -File .\elastic_log_clean.ps1 -InputPath .\local-diagnostics.zip
.\elastic_log_clean.ps1 -InputPath C:\logs\elasticsearch -VerboseLog
.\elastic_log_clean.ps1 -InputPath .\cluster.log -Config .\elastic_ip_name_map.yaml
.\elastic_log_clean.exe -v -i .\local-diagnostics.zip -o .\scrubbed-local-diagnostics.zip
```

`-v` / `-VerboseLog` prints every cleaned or omitted file. Stage lines (`extract:`, `process:`, `cleaned:`, `mappings:`) always print.

## Build

Requires Go 1.22+.

```bash
make
make test
```

Flags: `-c` `-i` `-o` `-w` `-r` `-v`. Do **not** share `report.yaml`.
