#!/usr/bin/env bash
set -euo pipefail
# Start the documented frontend port-forward first.
base_url=${1:-http://localhost:8080}
curl -fsS "$base_url/healthz"
curl -fsS "$base_url/config"
curl -fsS -H 'Content-Type: application/json' -d '{"text":"Instructor smoke test"}' "$base_url/api/notes"
curl -fsS "$base_url/api/notes" | python3 -c 'import json,sys; rows=json.load(sys.stdin); assert any(n["text"] == "Instructor smoke test" for n in rows); print("PASS: note persisted and listed")'
