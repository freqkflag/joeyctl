#!/usr/bin/env bash
set -euo pipefail

arch="$(uname -m)"
case "$arch" in
  x86_64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac

url="https://github.com/freqkflag/joeyctl/releases/latest/download/joeyctl-linux-${arch}"
tmp="$(mktemp)"
cleanup() { rm -f "$tmp"; }
trap cleanup EXIT

echo "Downloading ${url}..."
curl -fsSL "$url" -o "$tmp"

chmod +x "$tmp"
dest="/usr/local/bin/joeyctl"
dest_dir="$(dirname "$dest")"

install_cmd=(install -m 0755 "$tmp" "$dest")
if [ -w "$dest_dir" ]; then
  "${install_cmd[@]}"
else
  if command -v sudo >/dev/null 2>&1; then
    echo "Elevating with sudo to write to $dest_dir..."
    sudo "${install_cmd[@]}"
  else
    echo "Write access to $dest_dir is required; rerun as root or with sudo." >&2
    exit 1
  fi
fi

echo "joeyctl installed to $dest"
