#!/bin/sh
set -eu

forgejo_api=${FORGEJO_API_ROOT:-https://git.dannyhunn.com/api/v1}
forgejo_web=${FORGEJO_WEB_ROOT:-https://git.dannyhunn.com}
forgejo_remote=${GOTTH_FORGEJO_REMOTE:-origin}
expected=${EXPECT_STATE:-unpublished}
test_mode=${CLOSEOUT_TEST_MODE:-0}
tag=v1.0.0
if [ -z "${FORGEJO_TOKEN:-}" ]; then echo "FORGEJO_TOKEN is required" >&2; exit 2; fi
case "$expected" in unpublished|published) ;; *) echo "EXPECT_STATE must be unpublished or published" >&2; exit 2 ;; esac
if [ "$test_mode" != 1 ]; then
    [ "$forgejo_api" = https://git.dannyhunn.com/api/v1 ] &&
    [ "$forgejo_web" = https://git.dannyhunn.com ] || {
        echo "noncanonical Forgejo origin rejected" >&2
        exit 1
    }
fi
if remote_url=$(git remote get-url "$forgejo_remote" 2>/dev/null); then :; else remote_url=$forgejo_remote; fi
case "$remote_url" in
    forgejo:gotthboard/gotth-extension-webhook.git|ssh://forgejo@git.dannyhunn.com:56087/gotthboard/gotth-extension-webhook.git|https://git.dannyhunn.com/gotthboard/gotth-extension-webhook.git) ;;
    *) [ "$test_mode" = 1 ] || { echo "GOTTH_FORGEJO_REMOTE is not canonical Forgejo" >&2; exit 1; } ;;
esac

root=$(mktemp -d)
trap 'rm -rf "$root"' EXIT HUP INT TERM
umask 077
auth="$root/curl-auth.conf"
header_name=Authorization
auth_scheme=token
printf 'header = "%s: %s %s"\n' "$header_name" "$auth_scheme" "$FORGEJO_TOKEN" > "$auth"
repo="$root/repo.json"
code=$(curl --config "$auth" -sS -o "$repo" -w '%{http_code}' -H 'Accept: application/json' "$forgejo_api/repos/gotthboard/gotth-extension-webhook")
[ "$code" = 200 ] || { echo "Forgejo credential/repository check failed: HTTP $code" >&2; exit 1; }
python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); raise SystemExit(0 if d.get("full_name")=="gotthboard/gotth-extension-webhook" else 1)' "$repo"
refs=$(git ls-remote "$forgejo_remote" "refs/tags/$tag" "refs/tags/$tag^{}")
release="$root/release.json"
release_code=$(curl --config "$auth" -sS -o "$release" -w '%{http_code}' -H 'Accept: application/json' "$forgejo_api/repos/gotthboard/gotth-extension-webhook/releases/tags/$tag")
case "$release_code" in 200|404) ;; *) echo "Forgejo Release query failed: HTTP $release_code" >&2; exit 1 ;; esac
if [ "$expected" = unpublished ]; then
    [ -z "$refs" ] && [ "$release_code" = 404 ] || { echo "stable Forgejo publication exists; unpublished expectation is stale" >&2; exit 1; }
    echo 'forgejo_tag=absent forgejo_release=absent stable_closeout=unpublished'
    exit 0
fi

if [ -z "${EXPECTED_TAG_OBJECT:-}" ]; then echo "EXPECTED_TAG_OBJECT is required for published closeout" >&2; exit 2; fi
if [ -z "${EXPECTED_COMMIT:-}" ]; then echo "EXPECTED_COMMIT is required for published closeout" >&2; exit 2; fi
if [ -z "${BUNDLE_DIR:-}" ]; then echo "BUNDLE_DIR is required for published closeout" >&2; exit 2; fi
[ "$(git cat-file -t "refs/tags/$tag")" = tag ] || { echo "$tag is not annotated" >&2; exit 1; }
[ "$(git rev-parse "refs/tags/$tag")" = "$EXPECTED_TAG_OBJECT" ] || { echo "local tag object mismatch" >&2; exit 1; }
[ "$(git rev-parse "refs/tags/$tag^{}")" = "$EXPECTED_COMMIT" ] || { echo "local tag commit mismatch" >&2; exit 1; }
remote_object=$(printf '%s\n' "$refs" | awk -v r="refs/tags/$tag" '$2==r{print $1}')
remote_commit=$(printf '%s\n' "$refs" | awk -v r="refs/tags/$tag^{}" '$2==r{print $1}')
[ "$remote_object" = "$EXPECTED_TAG_OBJECT" ] && [ "$remote_commit" = "$EXPECTED_COMMIT" ] || { echo "Forgejo tag identity mismatch" >&2; exit 1; }
[ "$release_code" = 200 ] || { echo "stable Forgejo Release missing" >&2; exit 1; }
name=gotth-extension-webhook-1.0.0-linux-amd64
archive=$name.tar.gz
sbom=$name.spdx.json
provenance=$name.intoto.jsonl
for file in "$archive" "$sbom" "$provenance" SHA256SUMS; do [ -f "$BUNDLE_DIR/$file" ] || { echo "local release asset missing: $file" >&2; exit 1; }; done
(cd "$BUNDLE_DIR" && [ "$(wc -l < SHA256SUMS)" -eq 3 ] && sha256sum -c SHA256SUMS)
python3 - "$release" "$tag" "$archive" "$sbom" "$provenance" <<'PY'
import json,sys
d=json.load(open(sys.argv[1]))
expected=set(sys.argv[3:]+['SHA256SUMS'])
names=[a.get('name') for a in d.get('assets',[])]
ok=d.get('tag_name')==sys.argv[2] and d.get('draft') is False and d.get('prerelease') is False and len(names)==4 and set(names)==expected
raise SystemExit(0 if ok else 1)
PY
mkdir "$root/assets"
for file in "$archive" "$sbom" "$provenance" SHA256SUMS; do
    url=$(python3 - "$release" "$file" <<'PY'
import json,sys
d=json.load(open(sys.argv[1])); print({a['name']:a['browser_download_url'] for a in d['assets']}[sys.argv[2]])
PY
)
    case "$url" in "$forgejo_web/gotthboard/gotth-extension-webhook/releases/download/$tag/$file") ;; *) echo "unsafe Forgejo asset URL" >&2; exit 1 ;; esac
    curl --config "$auth" -fsSL -o "$root/assets/$file" "$url"
    cmp "$BUNDLE_DIR/$file" "$root/assets/$file"
done
(cd "$root/assets" && sha256sum -c SHA256SUMS)
echo 'forgejo_tag=verified forgejo_release=stable forgejo_assets=verified stable_closeout=complete'
