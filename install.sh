#!/bin/sh
# Install the latest prism release binary — no Go toolchain required.
#   curl -fsSL https://raw.githubusercontent.com/mcoder33/prism/main/install.sh | sh
# Override the install dir with PRISM_BINDIR (default /usr/local/bin).
set -eu

REPO="mcoder33/prism"
BINDIR="${PRISM_BINDIR:-/usr/local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*)
		echo "unsupported architecture: $arch" >&2
		exit 1
		;;
esac
case "$os" in
	linux | darwin) ;;
	*)
		echo "unsupported OS: $os — download the .zip from https://github.com/$REPO/releases on Windows" >&2
		exit 1
		;;
esac

tag=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
	grep '"tag_name"' | head -1 | cut -d'"' -f4)
if [ -z "$tag" ]; then
	echo "could not resolve the latest release tag" >&2
	exit 1
fi

asset="prism_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$tag/$asset"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $asset ($tag)…"
curl -fsSL "$url" -o "$tmp/$asset"
tar -xzf "$tmp/$asset" -C "$tmp"

if [ -w "$BINDIR" ]; then
	mv "$tmp/prism" "$BINDIR/prism"
else
	echo "Need elevated permissions to write $BINDIR (set PRISM_BINDIR to install elsewhere)…"
	sudo mv "$tmp/prism" "$BINDIR/prism"
fi
chmod +x "$BINDIR/prism"

echo "Installed prism to $BINDIR/prism"
"$BINDIR/prism" --version
