#!/usr/bin/env bash
# Build-if-needed wrapper, same idea as calling must-gather-clean after `make`.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="${ROOT}/bin/elastic_log_clean"
if [[ ! -x "${BIN}" ]]; then
  make -C "${ROOT}" build
fi
exec "${BIN}" "$@"
