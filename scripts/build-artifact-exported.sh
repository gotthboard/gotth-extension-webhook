#!/bin/sh
set -eu

if [ "$#" -ne 7 ] || [ "$1" != --immutable-export ]; then
    echo "internal immutable-export builder; use scripts/build-artifact.sh" >&2
    exit 2
fi
shift
version=$1
output=$2
commit=$3
tree=$4
epoch=$5
source_clean=$6
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
[ ! -e "$root/.git" ] && [ -f "$root/go.mod" ] || { echo "immutable source export unavailable" >&2; exit 1; }
printf '%s\n' "$commit" | grep -Eq '^[0-9a-f]{40}$' || { echo "invalid source commit identity" >&2; exit 1; }
printf '%s\n' "$tree" | grep -Eq '^[0-9a-f]{40}$' || { echo "invalid source tree identity" >&2; exit 1; }
printf '%s\n' "$epoch" | grep -Eq '^[0-9]+$' || { echo "invalid source epoch" >&2; exit 1; }
case $source_clean in 0|1) ;; *) echo "invalid source cleanliness state" >&2; exit 1 ;; esac

parent=$(dirname -- "$output")
base=$(basename -- "$output")
mkdir -p "$parent"
[ ! -e "$output" ] || { echo "release output path already exists: $output" >&2; exit 1; }
stage=$(mktemp -d "$parent/.${base}.stage.XXXXXX")
lock="$output.lock"
locked=0
cleanup() {
    rm -rf "$stage"
    if [ "$locked" = 1 ]; then rmdir "$lock" 2>/dev/null || true; fi
}
trap cleanup EXIT HUP INT TERM
mkdir "$lock" || { echo "release output publisher is already active" >&2; exit 1; }
locked=1
[ ! -e "$output" ] || { echo "release output path appeared during build: $output" >&2; exit 1; }

export GOENV=off GOWORK=off GOFLAGS= GOTOOLCHAIN=local GOEXPERIMENT= GOFIPS140=off
export CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1
export GOPROXY=off GOSUMDB=off GOPRIVATE= GONOPROXY= GONOSUMDB=
name="gotth-extension-webhook-$version-linux-amd64"
archive_name="$name.tar.gz"
sbom_name="$name.spdx.json"
provenance_name="$name.intoto.jsonl"
package_root="$stage/package"
package="$package_root/$name"
publish="$stage/publish"
mkdir -p "$package" "$publish"
(
cd "$root"
go mod verify
go list -mod=readonly -m -json all | python3 -c 'import json,sys
stream=json.JSONDecoder(); data=sys.stdin.read(); offset=0
while offset < len(data):
    while offset < len(data) and data[offset].isspace(): offset += 1
    if offset == len(data): break
    item,offset=stream.raw_decode(data,offset)
    if item.get("Replace") is not None: raise SystemExit("module replacement is forbidden in release builds")'
go build -mod=readonly -trimpath -buildvcs=false \
    -ldflags "-s -w -buildid= -X github.com/gotthboard/gotth-extension-webhook/internal/webhook.BuildVersion=$version" \
    -o "$package/gotth-extension-webhook" ./cmd/gotth-extension-webhook
go run -mod=readonly ./cmd/gotth-extension-webhook-manifest "$version" > "$package/manifest.json"
)
if [ "$("$package/gotth-extension-webhook" --version)" != "$version" ]; then
    echo "release binary version mismatch" >&2
    exit 1
fi
cp "$root/configuration-metadata.json" "$package/configuration-metadata.json"
cp "$root/LICENSE" "$package/LICENSE"
cp "$root/SECURITY.md" "$package/SECURITY.md"
python3 "$root/scripts/release-metadata.py" inventory \
    --root "$root" --version "$version" --binary "$package/gotth-extension-webhook" \
    --output "$package" --epoch "$epoch"
(
cd "$package"
find . -type f ! -name SHA256SUMS -print | LC_ALL=C sort | while IFS= read -r file; do sha256sum "$file"; done > SHA256SUMS
)
chmod 0755 "$package" "$package/gotth-extension-webhook" "$package/licenses"
chmod g-s "$package" "$package/licenses"
find "$package" -type f ! -name gotth-extension-webhook -exec chmod 0644 {} +
tar --format=ustar --sort=name --mtime="@$epoch" --owner=0 --group=0 --numeric-owner -C "$package_root" -cf - "$name" | gzip -n > "$publish/$archive_name"
cp "$package/SBOM.spdx.json" "$publish/$sbom_name"
clean_argument=
if [ "$source_clean" = 1 ]; then clean_argument=--clean; fi
python3 "$root/scripts/release-metadata.py" provenance \
    --root "$root" --version "$version" --binary "$package/gotth-extension-webhook" \
    --archive "$publish/$archive_name" --sbom "$publish/$sbom_name" \
    --source-manifest "$package/SOURCE_MANIFEST.sha256" \
    --commit "$commit" --tree "$tree" $clean_argument \
    --output "$publish/$provenance_name"
(cd "$publish" && sha256sum "$archive_name" "$sbom_name" "$provenance_name" > SHA256SUMS)
chmod 0644 "$publish/$archive_name" "$publish/$sbom_name" "$publish/$provenance_name" "$publish/SHA256SUMS"
for file in "$archive_name" "$sbom_name" "$provenance_name" SHA256SUMS; do sync -d "$publish/$file"; done
sync -d "$publish"
[ ! -e "$output" ] || { echo "release output path appeared before publication: $output" >&2; exit 1; }
mv -T -- "$publish" "$output"
sync -d "$parent"
archive_sha=$(sha256sum "$output/$archive_name" | cut -d' ' -f1)
rmdir "$lock"
locked=0
printf '%s  %s\n' "$archive_sha" "$output/$archive_name"
