#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
left=$(mktemp -d)
right=$(mktemp -d)
trap 'rm -rf "$left" "$right"' EXIT HUP INT TERM
epoch=$(git -C "$root" show -s --format=%ct HEAD)
SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" 1.0.0-alpha.2 "$left" >/dev/null
SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" 1.0.0-alpha.2 "$right" >/dev/null
cmp "$left/gotth-extension-webhook-1.0.0-alpha.2-linux-amd64.tar.gz" "$right/gotth-extension-webhook-1.0.0-alpha.2-linux-amd64.tar.gz"
archive="$left/gotth-extension-webhook-1.0.0-alpha.2-linux-amd64.tar.gz"
tar -tzf "$archive" | sort
test "$(tar -tvzf "$archive" | sed -n '1{s/[[:space:]].*//;p;}')" = drwxr-xr-x
extract=$(mktemp -d)
trap 'rm -rf "$left" "$right" "$extract"' EXIT HUP INT TERM
tar -xzf "$archive" -C "$extract"
name=gotth-extension-webhook-1.0.0-alpha.2-linux-amd64
test "$(stat -c %a "$extract/$name")" = 755
test "$("$extract/$name/gotth-extension-webhook" --version)" = 1.0.0-alpha.2
(cd "$extract/$name" && sha256sum -c SHA256SUMS)

before=$(sha256sum "$right/$name.tar.gz" | cut -d' ' -f1)
if SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" 1.0.0-alpha.2 "$right" >/dev/null 2>&1; then
    echo "existing release artifact was replaced" >&2
    exit 1
fi
after=$(sha256sum "$right/$name.tar.gz" | cut -d' ' -f1)
test "$before" = "$after"

if SOURCE_DATE_EPOCH=$((epoch + 1)) "$root/scripts/build-artifact.sh" 1.0.0-alpha.2 "$right/epoch" >/dev/null 2>&1; then
    echo "mismatched source epoch was accepted" >&2
    exit 1
fi
if SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" 1.0.0-alpha.2garbage "$right/version" >/dev/null 2>&1; then
    echo "invalid release version was accepted" >&2
    exit 1
fi
