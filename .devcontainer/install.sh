#!/usr/bin/env bash
set -e

# fns-cli
make install
fns-cli config import-extras
grep -qxF 'source ~/fns-cli.bash_profile' "$HOME/.bashrc" || printf '\nsource ~/fns-cli.bash_profile\n' >>"$HOME/.bashrc"

# goreleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Ensure gemini directory exists and is writable
mkdir -p /home/vscode/.gemini
sudo chown -R vscode:vscode /home/vscode/.gemini || true

if ! command -v gemini >/dev/null 2>&1; then
  echo "Installing Gemini CLI..."
  npm install -g @google/gemini-cli@latest
fi

if ! command -v shfmt &>/dev/null 2>&1; then
  echo "Installing shfmt..."
  go install mvdan.cc/sh/v3/cmd/shfmt@latest
fi
