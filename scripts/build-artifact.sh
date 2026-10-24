#!/bin/sh
set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <version> <new-output-directory>" >&2
    exit 2
fi
version=$1
output=$2
if ! printf '%s\n' "$version" | grep -Eq '^1\.0\.0$|^1\.0\.0-(alpha|beta)\.[1-9][0-9]*$'; then
    echo "invalid release version" >&2
    exit 2
fi
if [ -n "${GOTTH_RELEASE_EXPORT_ROOT:-}${GOTTH_RELEASE_COMMIT:-}${GOTTH_RELEASE_TREE:-}${GOTTH_RELEASE_EPOCH:-}${GOTTH_RELEASE_SOURCE_CLEAN:-}" ]; then
    echo "internal release-build environment is reserved" >&2
    exit 2
fi

checkout=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
commit=$(git -C "$checkout" rev-parse --verify 'HEAD^{commit}')
tree=$(git -C "$checkout" rev-parse "$commit^{tree}")
commit_epoch=$(git -C "$checkout" show -s --format=%ct "$commit")
epoch=${SOURCE_DATE_EPOCH:-$commit_epoch}
if [ "$epoch" != "$commit_epoch" ]; then
    echo "SOURCE_DATE_EPOCH must equal the source commit time" >&2
    exit 1
fi
require_clean=${GOTTH_EXTENSION_REQUIRE_CLEAN:-0}
case $require_clean in
    0|1) ;;
    *) echo "GOTTH_EXTENSION_REQUIRE_CLEAN must be 0 or 1" >&2; exit 2 ;;
esac
status=$(git -C "$checkout" status --porcelain=v1 --untracked-files=all)
source_clean=0
if [ -z "$status" ]; then
    source_clean=1
elif [ "$require_clean" = 1 ]; then
    echo "release artifact build requires a clean source tree" >&2
    exit 1
fi
export_stage=$(mktemp -d "${TMPDIR:-/tmp}/gotth-extension-webhook-export.XXXXXX")
cleanup_export() { rm -rf "$export_stage"; }
trap cleanup_export EXIT HUP INT TERM
mkdir "$export_stage/source"
git -C "$checkout" archive --format=tar "$commit" | tar -xf - -C "$export_stage/source"
"$export_stage/source/scripts/build-artifact-exported.sh" --immutable-export "$version" "$output" "$commit" "$tree" "$epoch" "$source_clean"
