#!/bin/sh
set -eu

forgejo_api=${FORGEJO_API_ROOT:-https://git.dannyhunn.com/api/v1}
forgejo_web=${FORGEJO_WEB_ROOT:-https://git.dannyhunn.com}
forgejo_remote=${GOTTH_FORGEJO_REMOTE:-origin}
github_api=${GITHUB_API_ROOT:-https://api.github.com}
github_web=${GITHUB_WEB_ROOT:-https://github.com}
test_mode=${CLOSEOUT_TEST_MODE:-0}
: "${FORGEJO_TOKEN:?FORGEJO_TOKEN is required}"
if [ "$test_mode" != 1 ]; then
 [ "$forgejo_api" = https://git.dannyhunn.com/api/v1 ] && [ "$forgejo_web" = https://git.dannyhunn.com ] && [ "$github_api" = https://api.github.com ] && [ "$github_web" = https://github.com ] || { echo "noncanonical HTTP origin rejected" >&2; exit 1; }
fi
tag=v1.0.0-alpha.1; expected=${EXPECT_STATE:-incomplete}
object_expected=a126253f637ee9cf8b60482640f5f99c128a2d24; commit_expected=3627c43213b82afa290372a6e2d9c36978a54758
archive=gotth-extension-webhook-1.0.0-alpha.1-linux-amd64.tar.gz; archive_sha=5cb6043ca200acfa67d4c6a85e0c1ba070c51dc550cacca7ce538021c8d9e83a; sums_sha=0480c3eabc44595fe8b2d011e53c4500010f8cf50a04293626b6c8471c44ca6b
case "$expected" in incomplete|complete) ;; *) echo "EXPECT_STATE must be incomplete or complete" >&2; exit 2 ;; esac
[ "$(git cat-file -t "refs/tags/$tag")" = tag ] || { echo "$tag is not annotated" >&2; exit 1; }
object=$(git rev-parse "refs/tags/$tag"); commit=$(git rev-parse "refs/tags/$tag^{}")
[ "$object" = "$object_expected" ] && [ "$commit" = "$commit_expected" ] || { echo "local identity mismatch" >&2; exit 1; }
if remote_url=$(git remote get-url "$forgejo_remote" 2>/dev/null); then :; else remote_url=$forgejo_remote; fi
case "$remote_url" in forgejo:gotthboard/gotth-extension-webhook.git|ssh://forgejo@git.dannyhunn.com:56087/gotthboard/gotth-extension-webhook.git|https://git.dannyhunn.com/gotthboard/gotth-extension-webhook.git) ;; *) echo "GOTTH_FORGEJO_REMOTE is not canonical Forgejo" >&2; exit 1 ;; esac
for remote in "$forgejo_remote" https://github.com/gotthboard/gotth-extension-webhook.git; do refs=$(git ls-remote "$remote" "refs/tags/$tag" "refs/tags/$tag^{}"); [ "$(printf '%s
' "$refs"|awk -v r="refs/tags/$tag" '$2==r{print $1}')" = "$object" ] && [ "$(printf '%s
' "$refs"|awk -v r="refs/tags/$tag^{}" '$2==r{print $1}')" = "$commit" ] || { echo "$remote tag identity mismatch" >&2; exit 1; }; done
root=$(mktemp -d); trap 'rm -rf "$root"' EXIT HUP INT TERM
probe="$root/forgejo-repo.json"; code=$(curl -sS -o "$probe" -w '%{http_code}' -H "Authorization: token $FORGEJO_TOKEN" -H 'Accept: application/json' "$forgejo_api/repos/gotthboard/gotth-extension-webhook")
[ "$code" = 200 ] || { echo "Forgejo credential/repository check failed: HTTP $code" >&2; exit 1; }
python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); sys.exit(0 if d.get("full_name")=="gotthboard/gotth-extension-webhook" else 1)' "$probe"
fj="$root/forgejo-release.json"; fjcode=$(curl -sS -o "$fj" -w '%{http_code}' -H "Authorization: token $FORGEJO_TOKEN" -H 'Accept: application/json' "$forgejo_api/repos/gotthboard/gotth-extension-webhook/releases/tags/$tag")
gh="$root/github-release.json"; ghcode=$(curl -sS -o "$gh" -w '%{http_code}' -H 'Accept: application/vnd.github+json' -H 'User-Agent: gotth-release-closeout' "$github_api/repos/gotthboard/gotth-extension-webhook/releases/tags/$tag")
if [ "$expected" = incomplete ]; then [ "$fjcode" = 404 ] && [ "$ghcode" = 404 ] || { echo "hosted Release state changed: forgejo=$fjcode github=$ghcode" >&2; exit 1; }; echo 'forgejo_release=absent github_release=not_published github_draft_state=unknown closeout=blocked'; exit 0; fi
[ "$fjcode" = 200 ] && [ "$ghcode" = 200 ] || { echo "hosted Release missing: forgejo=$fjcode github=$ghcode" >&2; exit 1; }
for body in "$fj" "$gh"; do python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); n=[a.get("name") for a in d.get("assets",[])]; ok=d.get("tag_name")==sys.argv[2] and d.get("draft") is False and d.get("prerelease") is True and len(n)==2 and set(n)=={sys.argv[3],"SHA256SUMS"}; sys.exit(0 if ok else 1)' "$body" "$tag" "$archive"; done
asset_url() { python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); print({a["name"]:a["browser_download_url"] for a in d["assets"]}[sys.argv[2]])' "$1" "$2"; }
fj_archive=$(asset_url "$fj" "$archive"); fj_sums=$(asset_url "$fj" SHA256SUMS); gh_archive=$(asset_url "$gh" "$archive"); gh_sums=$(asset_url "$gh" SHA256SUMS)
case "$fj_archive" in "$forgejo_web/gotthboard/gotth-extension-webhook/releases/download/$tag/$archive") ;; *) echo "unsafe Forgejo archive URL" >&2; exit 1 ;; esac; case "$fj_sums" in "$forgejo_web/gotthboard/gotth-extension-webhook/releases/download/$tag/SHA256SUMS") ;; *) echo "unsafe Forgejo checksum URL" >&2; exit 1 ;; esac
case "$gh_archive" in "$github_web/gotthboard/gotth-extension-webhook/releases/download/$tag/$archive") ;; *) echo "unsafe GitHub archive URL" >&2; exit 1 ;; esac; case "$gh_sums" in "$github_web/gotthboard/gotth-extension-webhook/releases/download/$tag/SHA256SUMS") ;; *) echo "unsafe GitHub checksum URL" >&2; exit 1 ;; esac
for host in forgejo github; do mkdir "$root/$host-assets"; done
curl -fsSL -H "Authorization: token $FORGEJO_TOKEN" -o "$root/forgejo-assets/$archive" "$fj_archive"; curl -fsSL -H "Authorization: token $FORGEJO_TOKEN" -o "$root/forgejo-assets/SHA256SUMS" "$fj_sums"
curl -fsSL -o "$root/github-assets/$archive" "$gh_archive"; curl -fsSL -o "$root/github-assets/SHA256SUMS" "$gh_sums"
for host in forgejo github; do printf '%s  %s
' "$archive_sha" "$root/$host-assets/$archive"|sha256sum -c -; printf '%s  %s
' "$sums_sha" "$root/$host-assets/SHA256SUMS"|sha256sum -c -; (cd "$root/$host-assets" && sha256sum -c SHA256SUMS); done
cmp "$root/forgejo-assets/$archive" "$root/github-assets/$archive"; cmp "$root/forgejo-assets/SHA256SUMS" "$root/github-assets/SHA256SUMS"
echo 'forgejo_release=present github_release=present prerelease=true assets=verified closeout=complete'
