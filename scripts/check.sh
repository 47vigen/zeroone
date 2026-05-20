#!/usr/bin/env bash
set -euo pipefail

# Pre-merge validation: tests, render xray config from the example,
# and (optionally) build the panel if dependencies are installed.
#
# Use `config/stack.local.json` if present (operators typically keep a
# real local copy there); otherwise fall back to the public example.

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
CONFIG="$ROOT/config/stack.local.json"
[ -f "$CONFIG" ] || CONFIG="$ROOT/config/stack.example.json"

go vet ./...
go test ./...
go run ./cmd/xray-stackd -config "$CONFIG" -print-xray >/tmp/xray-stackd-generated.json
python3 -m json.tool /tmp/xray-stackd-generated.json >/dev/null

if [ -d "$ROOT/web/app/node_modules" ]; then
  npm --prefix "$ROOT/web/app" run build >/dev/null
fi
