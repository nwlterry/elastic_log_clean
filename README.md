# elastic_log_clean

Go rewrite of the Elastic diagnostic / cluster / server log cleaner, using the same omit + obfuscate + `report.yaml` model as [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean).

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

## Download

Pre-built binaries are on the [Releases](https://github.com/nwlterry/elastic_log_clean/releases) page (`v1.1.0`).

| Archive | Use on |
|---------|--------|
| `elastic_log_clean_1.1.0_linux_amd64.tar.gz` | RHEL / most servers |
| `elastic_log_clean_1.1.0_linux_arm64.tar.gz` | Linux aarch64 |
| `elastic_log_clean_1.1.0_darwin_amd64.tar.gz` | Intel macOS |
| `elastic_log_clean_1.1.0_darwin_arm64.tar.gz` | Apple Silicon |
| `elastic_log_clean_1.1.0_windows_amd64.zip` | Windows |

```bash
curl -fsSL -O https://github.com/nwlterry/elastic_log_clean/releases/download/v1.1.0/elastic_log_clean_1.1.0_linux_amd64.tar.gz
tar -xzf elastic_log_clean_1.1.0_linux_amd64.tar.gz
chmod +x elastic_log_clean
./elastic_log_clean --version
./elastic_log_clean -i diagnostics.zip -o scrubbed-diagnostics.zip
```

Each archive includes the binary plus `elastic_default.yaml` and `elastic_ip_name_map.yaml`.

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

## Config: IP only vs IP + custom names

Passing `-c` **replaces** the built-in config. It does not merge. A file that only contains `type: IP` will rewrite IPs and skip MAC / email / secrets / PEM / omit.

Default (`type: IP`, `replacementType: Consistent`) assigns tokens:

`10.99.1.8` → `x-ipv4-0000000001-x` (same source IP always gets the same token). Loopback `127.0.0.1` / `0.0.0.0` / `::1` is left alone.

The `IP` rule has no name map. Pin specific IPs or hostnames with `Keywords`. Keywords run **before** IP, so a mapped address is not tokenized again.

Copy [config/elastic_ip_name_map.yaml](config/elastic_ip_name_map.yaml) and edit the `replacement` maps.

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
