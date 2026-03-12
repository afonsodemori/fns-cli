#!/usr/bin/env bash
set -e

# fns-cli
brew remove fns-cli || true
make dev
fns-cli config import-extras
grep -qxF 'source ~/fns-cli.bash_profile' "$HOME/.bashrc" || printf '\nsource ~/fns-cli.bash_profile\n' >>"$HOME/.bashrc"

# Ensure gemini directory is writable
sudo chown -R vscode:vscode /home/vscode/.gemini || true
