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
echo "node 10.20.30.41 password=secret" | ./bin/elastic_log_clean --stdin
./bin/elastic_log_clean --print-default-config
```

Flags match must-gather-clean: `-c` `-i` `-o` `-w` `-r`.

Do **not** share `report.yaml`.

## Layout

```
cmd/elastic_log_clean/main.go
pkg/clean/
config/elastic_default.yaml
examples/sample-elasticsearch.log
scripts/elastic_log_clean.sh
Makefile
```

See [GROUP.md](GROUP.md).
