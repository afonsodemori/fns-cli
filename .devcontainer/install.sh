#!/usr/bin/env bash
set -e

# fns-cli -- defensive, with "|| true" in case dev repo is failling state
make dev || true
fns-cli config import-extras || true
grep -qxF 'source ~/fns-cli.bash_profile' "$HOME/.bashrc" || printf '\nsource ~/fns-cli.bash_profile\n' >>"$HOME/.bashrc"

# Ensure gemini directory is writable
sudo chown -R vscode:vscode /home/vscode/.gemini || true
