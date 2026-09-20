#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
left=$(mktemp -d)
right=$(mktemp -d)
trap 'rm -rf "$left" "$right"' EXIT HUP INT TERM
SOURCE_DATE_EPOCH=1790000000 "$root/scripts/build-artifact.sh" 1.0.0-alpha.1 "$left" >/dev/null
SOURCE_DATE_EPOCH=1790000000 "$root/scripts/build-artifact.sh" 1.0.0-alpha.1 "$right" >/dev/null
cmp "$left/gotth-extension-webhook-1.0.0-alpha.1-linux-amd64.tar.gz" "$right/gotth-extension-webhook-1.0.0-alpha.1-linux-amd64.tar.gz"
tar -tzf "$left/gotth-extension-webhook-1.0.0-alpha.1-linux-amd64.tar.gz" | sort
