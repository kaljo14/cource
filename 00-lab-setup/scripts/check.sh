#!/usr/bin/env bash
set -euo pipefail
for tool in docker kind kubectl go git; do command -v "$tool"; done
docker info >/dev/null
