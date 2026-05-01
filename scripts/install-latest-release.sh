#!/usr/bin/env bash
set -euo pipefail

repo="nicolaiort/shikai"
install_dir="${INSTALL_DIR:-${1:-$HOME/bin}}"

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  linux|darwin) ;;
  *)
    echo "this installer only supports linux and macos" >&2
    exit 1
    ;;
esac

binary_name="shikai"
asset_name="shikai-${os}-${arch}"
destination="${install_dir}/${binary_name}"

api_url="https://api.github.com/repos/${repo}/releases/latest"
download_url="$(
  curl -fsSL -H 'Accept: application/vnd.github+json' -H 'User-Agent: shikai-installer' "$api_url" |
    awk -v asset_name="$asset_name" '
      $0 ~ /"name":/ && $0 ~ asset_name { found=1 }
      found && $0 ~ /"browser_download_url":/ {
        match($0, /"browser_download_url": *"([^"]+)"/, m)
        if (m[1] != "") {
          print m[1]
          exit
        }
      }
    '
)"

if [[ -z "$download_url" ]]; then
  echo "could not find asset '$asset_name' in latest release" >&2
  exit 1
fi

mkdir -p "$install_dir"
curl -fsSL -H 'Accept: application/octet-stream' -o "$destination" "$download_url"
chmod +x "$destination"

echo "installed $binary_name to $destination"
