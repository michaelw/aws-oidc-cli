#!/bin/sh
# Installs pre-commit and sets up git hooks
set -e
if ! command -v pre-commit >/dev/null 2>&1; then
  echo "Installing pre-commit..."
  pip install pre-commit
fi
pre-commit install
