#!/usr/bin/env bash
set -euo pipefail
kind create cluster --name course --config 00-lab-setup/kind-config.yaml
kubectl --context kind-course wait --for=condition=Ready nodes --all --timeout=180s
