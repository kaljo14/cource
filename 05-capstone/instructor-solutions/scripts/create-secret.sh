#!/usr/bin/env bash
set -euo pipefail
if [ "$(kubectl config current-context)" != "kind-capstone" ]; then
  echo "Select kind-capstone before running this script." >&2
  exit 1
fi
# Do not replace an existing database password; PostgreSQL initializes it only once.
if kubectl -n capstone get secret db-credentials >/dev/null 2>&1; then
  echo "Existing db-credentials retained."
  exit 0
fi
password_file=$(mktemp)
trap 'rm -f "$password_file"' EXIT
chmod 600 "$password_file"
openssl rand -hex 24 | tr -d '\n' > "$password_file"
kubectl -n capstone create secret generic db-credentials --from-file=password="$password_file"
