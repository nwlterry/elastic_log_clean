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

## Matcher guards (v1.2.1)

- Java / Maven coords such as `io.netty.transport@4.1.135.Final` stay as-is.
- Clock values such as `07:29:19` are not treated as IPv6.
- Real addresses (`sre@corp.example.com`, `2001:db8::1`) are still tokenized.

Do **not** share `report.yaml`.
