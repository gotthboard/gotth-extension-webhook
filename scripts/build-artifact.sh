#!/bin/sh
set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <version> <output-directory>" >&2
    exit 2
fi

version=$1
output=$2
case "$version" in
    1.0.0-alpha.[1-9]*|1.0.0-beta.[1-9]*|1.0.0) ;;
    *) echo "invalid release version" >&2; exit 2 ;;
esac

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
epoch=${SOURCE_DATE_EPOCH:-0}
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
name="gotth-extension-webhook-$version-linux-amd64"
mkdir -p "$stage/$name" "$output"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags "-s -w -buildid= -X github.com/gotthboard/gotth-extension-webhook/internal/webhook.BuildVersion=$version" \
    -o "$stage/$name/gotth-extension-webhook" "$root/cmd/gotth-extension-webhook"
go run "$root/cmd/gotth-extension-webhook-manifest" "$version" > "$stage/$name/manifest.json"
cp "$root/configuration-metadata.json" "$stage/$name/configuration-metadata.json"
cp "$root/LICENSE" "$stage/$name/LICENSE"
(cd "$stage/$name" && sha256sum LICENSE configuration-metadata.json gotth-extension-webhook manifest.json > SHA256SUMS)
chmod 0755 "$stage/$name/gotth-extension-webhook"
chmod 0644 "$stage/$name/LICENSE" "$stage/$name/configuration-metadata.json" "$stage/$name/manifest.json" "$stage/$name/SHA256SUMS"
archive="$output/$name.tar.gz"
tar --sort=name --mtime="@$epoch" --owner=0 --group=0 --numeric-owner -C "$stage" -cf - "$name" | gzip -n > "$archive"
sha256sum "$archive"
