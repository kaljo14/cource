#!/usr/bin/env bash
set -euo pipefail

REPO_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$REPO_DIR"
: "${STACKIT_REPO_URL:?Set STACKIT_REPO_URL to the SSH or HTTPS clone URL of your existing STACKIT Git repository}"

# Use SSH or a Git credential helper; never embed a token in the URL.
case "$STACKIT_REPO_URL" in
  https://*)
    if [[ "${STACKIT_REPO_URL#https://}" == *@* ]]; then
      printf 'Use an HTTPS URL without embedded credentials.\n' >&2
      exit 1
    fi
    ;;
  git@*:*|ssh://*) ;;
  *) printf 'Use an SSH or HTTPS repository clone URL.\n' >&2; exit 1 ;;
esac

if [ ! -e .git ]; then
  git init -b main
fi
BRANCH=$(git branch --show-current)
case "$BRANCH" in
  main|kaloyan-app) ;;
  *) printf 'Switch to main or kaloyan-app before running this script.\n' >&2; exit 1 ;;
esac
if git remote get-url origin >/dev/null 2>&1; then
  if [ "$(git remote get-url origin)" != "$STACKIT_REPO_URL" ]; then
    printf 'origin differs from STACKIT_REPO_URL; check your remote first.\n' >&2
    exit 1
  fi
else
  git remote add origin "$STACKIT_REPO_URL"
fi

# The caller reviews and commits changes explicitly before publication.
if ! git rev-parse --verify HEAD >/dev/null 2>&1; then
  printf 'Review, stage, and commit the course files before pushing.\n' >&2
  exit 1
fi
if [ -n "$(git status --porcelain)" ]; then
  printf 'Commit or stash your changes before pushing so CI matches your working tree.\n' >&2
  exit 1
fi
git push -u origin "$BRANCH"
printf 'Open the repository Actions tab to follow course-ci.yaml and capstone-ci.yaml.\n'
