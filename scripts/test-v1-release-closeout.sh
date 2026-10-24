#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d)
server_pid=
cleanup() {
    if [ -n "$server_pid" ]; then kill "$server_pid" 2>/dev/null || true; wait "$server_pid" 2>/dev/null || true; fi
    rm -rf "$tmp"
}
trap cleanup EXIT HUP INT TERM
git clone -q "$root" "$tmp/repo"
git -C "$tmp/repo" -c user.name=Closeout -c user.email=closeout@localhost tag -a v1.0.0 -m v1.0.0
git clone -q --bare "$tmp/repo" "$tmp/remote.git"
object=$(git -C "$tmp/repo" rev-parse refs/tags/v1.0.0)
commit=$(git -C "$tmp/repo" rev-parse 'refs/tags/v1.0.0^{}')
name=gotth-extension-webhook-1.0.0-linux-amd64
bundle=$tmp/bundle
mkdir "$bundle"
printf archive > "$bundle/$name.tar.gz"
printf sbom > "$bundle/$name.spdx.json"
printf provenance > "$bundle/$name.intoto.jsonl"
(cd "$bundle" && sha256sum "$name.tar.gz" "$name.spdx.json" "$name.intoto.jsonl" > SHA256SUMS)
cat > "$tmp/server.py" <<'PY'
import http.server,json,os,socketserver,sys
bundle,port_file,expected_token=sys.argv[1:]
class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.headers.get("Authorization") != "token "+expected_token:
            self.send_error(401)
            return
        if self.path.endswith('/repos/gotthboard/gotth-extension-webhook'):
            self.reply(200,{'full_name':'gotthboard/gotth-extension-webhook'})
            return
        if self.path.endswith('/releases/tags/v1.0.0'):
            base=f'http://127.0.0.1:{self.server.server_address[1]}/gotthboard/gotth-extension-webhook/releases/download/v1.0.0'
            names=os.listdir(bundle)
            self.reply(200,{'tag_name':'v1.0.0','draft':False,'prerelease':False,'assets':[{'name':n,'browser_download_url':f'{base}/{n}'} for n in names]})
            return
        marker='/releases/download/v1.0.0/'
        if marker in self.path:
            name=self.path.split(marker,1)[1]
            path=os.path.join(bundle,name)
            if os.path.isfile(path):
                data=open(path,'rb').read(); self.send_response(200); self.send_header('Content-Length',str(len(data))); self.end_headers(); self.wfile.write(data); return
        self.send_error(404)
    def reply(self,code,value):
        data=json.dumps(value).encode(); self.send_response(code); self.send_header('Content-Type','application/json'); self.send_header('Content-Length',str(len(data))); self.end_headers(); self.wfile.write(data)
    def log_message(self,*args): pass
with socketserver.TCPServer(('127.0.0.1',0),Handler) as server:
    open(port_file,'w').write(str(server.server_address[1]))
    server.serve_forever()
PY
python3 "$tmp/server.py" "$bundle" "$tmp/port" test-token &
server_pid=$!
i=0
while [ ! -s "$tmp/port" ]; do
    kill -0 "$server_pid"
    i=$((i+1)); [ "$i" -lt 100 ]
    sleep 0.01
done
port=$(cat "$tmp/port")
(
cd "$tmp/repo"
FORGEJO_TOKEN=test-token CLOSEOUT_TEST_MODE=1 EXPECT_STATE=published \
FORGEJO_API_ROOT="http://127.0.0.1:$port/api/v1" FORGEJO_WEB_ROOT="http://127.0.0.1:$port" \
GOTTH_FORGEJO_REMOTE="$tmp/remote.git" EXPECTED_TAG_OBJECT="$object" EXPECTED_COMMIT="$commit" BUNDLE_DIR="$bundle" \
"$root/scripts/check-v1-release-closeout.sh"
)
if (
cd "$tmp/repo" && FORGEJO_TOKEN=test-token CLOSEOUT_TEST_MODE=1 EXPECT_STATE=published \
FORGEJO_API_ROOT="http://127.0.0.1:$port/api/v1" FORGEJO_WEB_ROOT="http://127.0.0.1:$port" \
GOTTH_FORGEJO_REMOTE="$tmp/remote.git" EXPECTED_TAG_OBJECT="$object" EXPECTED_COMMIT="$commit" BUNDLE_DIR="$bundle" \
sh -c 'printf tamper >> "$BUNDLE_DIR/'"$name"'.tar.gz"; exec "$0"' "$root/scripts/check-v1-release-closeout.sh"
); then
    echo 'tampered stable asset accepted' >&2
    exit 1
fi
printf '%s\n' 'stable Forgejo closeout verifies exact tag, metadata, asset set, and bytes'
