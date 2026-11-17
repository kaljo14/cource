#!/usr/bin/env bash
set -euo pipefail

MODULE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
IMAGE=${IMAGE:-course-app:v1}
VERSION=${VERSION:-v1}

printf '\n==> Test the Go application\n'
(cd "$MODULE_DIR/app" && go test -race ./...)

printf '\n==> Build the application and container image: %s\n' "$IMAGE"
docker build --build-arg "VERSION=$VERSION" --tag "$IMAGE" "$MODULE_DIR/app"
