#!/usr/bin/env bash
set -euo pipefail
# Deletes all cluster workloads and node-local data.
kind delete cluster --name course
