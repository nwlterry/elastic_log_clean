# elastic_log_clean

Go rewrite of the Elastic diagnostic / cluster / server log cleaner, using the same omit + obfuscate + `report.yaml` model as [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean).

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

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
