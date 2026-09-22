# elastic_log_clean

Obfuscate sensitive data in **Elastic Support Diagnostics**, extracted diagnostic trees, cluster logs, node logs, Kibana/Fleet logs, and `elasticsearch.yml` dumps.

Logic follows [openshift/must-gather-clean](https://github.com/openshift/must-gather-clean):

- `omit` — drop files (keystores, PKCS12, `users` / `users_roles`)
- `obfuscate` — rewrite contents (and optionally paths) with consistent or static tokens
- `report.yaml` — mapping of original to token (**keep private**)

Sibling: [kafka_log_clean](https://github.com/nwlterry/kafka_log_clean).

## Default redactions

- IPv4/IPv6 (not `127.0.0.1`, `0.0.0.0`, `::1`) → `x-ipv4-0000000001-x` (stable in one run)
- MAC → `x-mac-0000000001-x`
- Email → `x-email-0000000001-x`
- `Authorization` / `ApiKey` / Bearer / Basic
- `password`, `bootstrap.password`, `bind_password`, `cloud.auth`, keystore passwords, service tokens
- `user:pass@host` URL passwords
- PEM private-key blocks
- Omit `*.p12`, `*.jks`, `elasticsearch.keystore`, `users`, `users_roles`

Version strings such as `8.18.4` are left alone. This is not a guarantee of complete sanitization. Do not attach `report.yaml` to a case.

## Requirements

Python 3.9+ and PyYAML:

```bash
pip install -r requirements.txt
```

## Usage

```bash
python3 elastic_log_clean.py -i local-diagnostics-20260922.zip -o scrubbed-local-diagnostics-20260922.zip
python3 elastic_log_clean.py -i /var/log/elasticsearch -o /tmp/es-logs-cleaned
python3 elastic_log_clean.py -i prod_server.json -o prod_server.cleaned.json
echo "node 10.20.30.41 password=secret" | python3 elastic_log_clean.py --stdin
python3 elastic_log_clean.py -c config/elastic_default.yaml -i diagnostics.zip -o scrubbed-diagnostics.zip -r /tmp/elastic-clean-report.yaml -w 8
bash scripts/elastic_log_clean.sh --version
python3 tests/test_clean.py
```

Flags: `-c` config, `-i` input file/dir/zip, `-o` output, `-w` workers, `-r` report, `--stdin`, `--print-default-config`.

## Config schema

```yaml
config:
  omit:
    - type: File
      pattern: "**/*.p12"
  obfuscate:
    - type: IP
      replacementType: Consistent   # or Static
      target: All                   # FileContents | FilePath | All
    - type: MAC
      replacementType: Consistent
    - type: Email
      replacementType: Consistent
    - type: Domain
      replacementType: Consistent
      domainNames: [corp.example.com]
    - type: Keywords
      target: FileContents
      replacement:
        es-prod-01: node-a
    - type: Regex
      regex: '(?i)(cluster\.uuid\s*[:=]\s*)(\S+)'
      replacement: '\1x-redacted-uuid-x'
    - type: Secrets
      enabled: true
    - type: PEM
      enabled: true
```

See [GROUP.md](GROUP.md). Catalog: https://github.com/nwlterry/nwlterry
