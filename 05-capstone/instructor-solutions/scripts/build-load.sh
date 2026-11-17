#!/usr/bin/env bash
set -euo pipefail
if [ "$(kubectl config current-context)" != "kind-capstone" ]; then
  echo "Select kind-capstone before running this script." >&2
  exit 1
fi
# Run from this instructor-solutions directory. These are visible, ordinary commands.
docker build -f api/Dockerfile --build-arg VERSION=v1 -t campus-api:v1 api
docker build -f frontend/Dockerfile -t campus-frontend:v1 frontend
kind load docker-image campus-api:v1 campus-frontend:v1 --name capstone
