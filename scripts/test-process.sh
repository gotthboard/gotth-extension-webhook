#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d)
pid=
cleanup() {
    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
        kill -TERM "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
    fi
    rm -rf "$tmp"
}
trap cleanup EXIT HUP INT TERM
version=1.0.0
binary=$tmp/gotth-extension-webhook
CGO_ENABLED=0 go build -trimpath -buildvcs=false \
    -ldflags "-X github.com/gotthboard/gotth-extension-webhook/internal/webhook.BuildVersion=$version" \
    -o "$binary" "$root/cmd/gotth-extension-webhook"
[ "$("$binary" --version)" = "$version" ]
if "$binary" >"$tmp/missing.out" 2>"$tmp/missing.err"; then
    echo 'process accepted missing runtime configuration' >&2
    exit 1
fi
[ ! -s "$tmp/missing.out" ]
[ "$(cat "$tmp/missing.err")" = 'gotth-extension-webhook: startup failed' ]

mkdir -m 0700 "$tmp/secrets"
printf '%s' "$(printf 'k%.0s' $(seq 1 32))" >"$tmp/secrets/webhook.hmac-key"
printf '%s' "$(printf 't%.0s' $(seq 1 32))" >"$tmp/token"
printf '%s\n' '{"schema":"gotth.extension.webhook.config.v1","endpoint":"https://receiver.example.test/hook","timeout_seconds":5}' >"$tmp/config.json"
manifest=$tmp/manifest.json
(cd "$root" && go run ./cmd/gotth-extension-webhook-manifest "$version") >"$manifest"
manifest_sha=$(sha256sum "$manifest" | cut -d' ' -f1)
python3 - "$manifest_sha" "$tmp/binding.json" <<'PY'
import hashlib,json,sys
manifest_sha,output=sys.argv[1:]
instance="00000000-0000-4000-8000-000000000019"
capabilities=sorted(["notification.alert.send","notification.alert.sink","notification.delivery.status"])
interfaces=[{"name":"gotth.mail.interface.notification","major":1,"minor":0}]
grant={"schema":"gotth.extensions.grant.v1","instance_id":instance,"extension_id":"gotth.mail.notification.webhook","manifest_sha256":manifest_sha,"capabilities":capabilities,"interfaces":interfaces,"secrets":["webhook.hmac-key"]}
canonical_grant=json.dumps(grant,separators=(",",":"))
grant_sha=hashlib.sha256(canonical_grant.encode()).hexdigest()
session={"schema":"gotth.extensions.session.v1","extension_id":"gotth.mail.notification.webhook","extension_version":"1.0.0","instance_id":instance,"control":{"name":"gotth.extensions.control","major":1,"minor":0},"interfaces":interfaces,"capabilities":capabilities,"secrets":["webhook.hmac-key"],"manifest_sha256":manifest_sha,"grant_sha256":grant_sha}
session_sha=hashlib.sha256(json.dumps(session,separators=(",",":")).encode()).hexdigest()
host={"protocols":[{"name":"gotth.extensions.control","major":1,"min_minor":0,"max_minor":0}],"interfaces":[{"name":"gotth.mail.interface.notification","major":1,"min_minor":0,"max_minor":0}]}
binding={"schema":"gotth.extension.runtime-binding.v1","instance_id":instance,"manifest_sha256":manifest_sha,"grant_sha256":grant_sha,"session_sha256":session_sha,"grant":grant,"host":host}
open(output,"w",encoding="utf-8").write(json.dumps(binding,separators=(",",":"))+"\n")
PY
chmod 0600 "$tmp/config.json" "$tmp/binding.json" "$tmp/token" "$tmp/secrets/webhook.hmac-key"
socket=$tmp/extension.sock
GOTTH_EXTENSION_CONFIG_FILE=$tmp/config.json \
GOTTH_EXTENSION_BINDING_FILE=$tmp/binding.json \
GOTTH_EXTENSION_SECRET_DIR=$tmp/secrets \
GOTTH_EXTENSION_SERVICE_TOKEN_FILE=$tmp/token \
GOTTH_EXTENSION_SOCKET=$socket \
"$binary" >"$tmp/process.out" 2>"$tmp/process.err" &
pid=$!
i=0
while [ ! -S "$socket" ]; do
    kill -0 "$pid" 2>/dev/null || { cat "$tmp/process.err" >&2; exit 1; }
    i=$((i + 1))
    [ "$i" -lt 100 ] || { echo 'process socket startup timed out' >&2; exit 1; }
    sleep 0.02
done
[ "$(stat -c %a "$socket")" = 600 ]
kill -TERM "$pid"
wait "$pid"
pid=
[ ! -e "$socket" ]
[ ! -s "$tmp/process.out" ]
[ ! -s "$tmp/process.err" ]
printf '%s\n' 'version, fail-closed startup, socket mode, and graceful process lifecycle verified'
