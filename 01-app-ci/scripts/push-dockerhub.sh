#!/usr/bin/env bash
set -euo pipefail

IMAGE=${IMAGE:-course-app:v1}
IMAGE_TAG=${IMAGE_TAG:-v1}
if [ -z "${DOCKERHUB_USER:-}" ]; then
  read -r -p 'Docker Hub username: ' DOCKERHUB_USER
fi
if [ -z "${DOCKERHUB_TOKEN:-}" ]; then
  read -r -s -p 'Docker Hub token (read/write): ' DOCKERHUB_TOKEN
  printf '\n'
fi
: "${DOCKERHUB_USER:?Docker Hub username is required}"
: "${DOCKERHUB_TOKEN:?Docker Hub token is required}"
TARGET_IMAGE="$DOCKERHUB_USER/course-app:$IMAGE_TAG"

printf '\n==> Log in to Docker Hub\n'
printf '%s' "$DOCKERHUB_TOKEN" | docker login --username "$DOCKERHUB_USER" --password-stdin
unset DOCKERHUB_TOKEN

printf '\n==> Tag %s as %s\n' "$IMAGE" "$TARGET_IMAGE"
docker tag "$IMAGE" "$TARGET_IMAGE"

printf '\n==> Push %s\n' "$TARGET_IMAGE"
docker push "$TARGET_IMAGE"
