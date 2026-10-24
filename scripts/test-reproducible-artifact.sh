#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
left_root=$(mktemp -d)
right_root=$(mktemp -d)
extract=$(mktemp -d)
clone_root=$(mktemp -d)
replace_root=$(mktemp -d)
relative_root=$(mktemp -d)
left=$left_root/bundle
right=$right_root/bundle
trap 'rm -rf "$left_root" "$right_root" "$extract" "$clone_root" "$replace_root" "$relative_root"' EXIT HUP INT TERM
export GOENV=off GOWORK=off GOFLAGS= GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOPRIVATE= GONOPROXY= GONOSUMDB= GOAMD64=v1
epoch=$(git -C "$root" show -s --format=%ct HEAD)
version=1.0.0
name=gotth-extension-webhook-$version-linux-amd64
archive=$name.tar.gz
sbom=$name.spdx.json
provenance=$name.intoto.jsonl
SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" "$version" "$left" >/dev/null
real_go=$(command -v go)
mkdir "$right_root/bin"
cat > "$right_root/bin/go" <<'SH'
#!/bin/sh
set -eu
[ "$GOENV" = off ] && [ "$GOWORK" = off ] && [ -z "$GOFLAGS" ] && [ "$GOTOOLCHAIN" = local ]
[ -z "$GOEXPERIMENT" ] && [ "$GOFIPS140" = off ]
[ "$GOPROXY" = off ] && [ "$GOSUMDB" = off ] && [ "$GOAMD64" = v1 ] && [ "$CGO_ENABLED" = 0 ]
printf 'sanitized\n' >> "$GO_AUDIT"
exec "$REAL_GO" "$@"
SH
chmod 0755 "$right_root/bin/go"
env PATH="$right_root/bin:$PATH" REAL_GO="$real_go" GO_AUDIT="$right_root/go-audit" SOURCE_DATE_EPOCH="$epoch" GOENV="$right_root/hostile-goenv" GOWORK="$right_root/hostile.work" GOFLAGS="-mod=mod -tags=hostile_ambient_flag" GOTOOLCHAIN=auto GOEXPERIMENT=fieldtrack GOFIPS140=latest GOPROXY=https://invalid.example.test GOSUMDB=sum.golang.org GOAMD64=v4 CGO_ENABLED=1 "$root/scripts/build-artifact.sh" "$version" "$right" >/dev/null
test -s "$right_root/go-audit"
if GOTTH_RELEASE_EXPORT_ROOT="$root" "$root/scripts/build-artifact.sh" "$version" "$right_root/reserved" >"$right_root/reserved.log" 2>&1; then
    echo 'reserved internal release environment was accepted' >&2
    exit 1
fi
grep -Fq 'internal release-build environment is reserved' "$right_root/reserved.log"
for file in "$archive" "$sbom" "$provenance" SHA256SUMS; do cmp "$left/$file" "$right/$file"; done
mkdir "$relative_root/caller"
(
cd "$relative_root/caller"
SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" "$version" ../bundle >/dev/null
)
for file in "$archive" "$sbom" "$provenance" SHA256SUMS; do cmp "$left/$file" "$relative_root/bundle/$file"; done
(cd "$left" && sha256sum -c SHA256SUMS)
tar -xzf "$left/$archive" -C "$extract"
test "$(stat -c %a "$extract/$name")" = 755
test "$(stat -c %a "$extract/$name/gotth-extension-webhook")" = 755
test "$("$extract/$name/gotth-extension-webhook" --version)" = "$version"
(cd "$extract/$name" && sha256sum -c SHA256SUMS)
cmp "$left/$sbom" "$extract/$name/SBOM.spdx.json"
python3 - "$left/$archive" "$left/$sbom" "$left/$provenance" "$root" <<'PY'
import base64,hashlib,json,subprocess,sys
archive,sbom,provenance,root=sys.argv[1:]
with open(sbom, encoding='utf-8') as handle:
    s=json.load(handle)
assert s['spdxVersion']=='SPDX-2.3'
assert s['name']=='gotth-extension-webhook-1.0.0-linux-amd64'
packages={(p['name'],p['versionInfo'],p['licenseDeclared']) for p in s['packages']}
expected={
 ('gotth-extension-webhook','1.0.0','MIT'),
 ('github.com/gotthboard/gotth-extensions','v1.0.0','MIT'),
 ('github.com/gotthboard/gotth-sdk','v1.0.0','MIT'),
 ('golang.org/x/net','v0.60.0','BSD-3-Clause'),
 ('golang.org/x/sys','v0.48.0','BSD-3-Clause'),
 ('golang.org/x/text','v0.42.0','BSD-3-Clause'),
 ('google.golang.org/genproto/googleapis/rpc','v0.0.0-20260526163538-3dc84a4a5aaa','Apache-2.0'),
 ('google.golang.org/grpc','v1.83.2','Apache-2.0'),
 ('google.golang.org/protobuf','v1.36.11','BSD-3-Clause'),
}
assert packages==expected
expected_zip={
 'github.com/gotthboard/gotth-extensions':'f5caba0cb6db41f0b53887316d57450ce7482a78921a12d770d884aaca6ff4e3',
 'github.com/gotthboard/gotth-sdk':'c975e47838b9882ca1d056fa2882ab220eb3e6aecb8f37ff1feb6776cdc37102',
 'golang.org/x/net':'052e0453db47a96ae9567503e7c95ea925ab52682cb1cec31250541f7160107e',
 'golang.org/x/sys':'3a47e06431ba1e24d5b2aca3a7ce31fd13faff34df3e5d8e63f105ff9d0c1786',
 'golang.org/x/text':'a7b64e003056b6470303f408202098d8f3714a115f23091b8cac85edeb265476',
 'google.golang.org/genproto/googleapis/rpc':'55f932b106d6f4f8cb76ae869e44f5d6d8400965fa441ee3532d09c96d64aac1',
 'google.golang.org/grpc':'84292905a694e17862d15abc123a6f7f7efa9e1e3a77e80120982ff95b2605ea',
 'google.golang.org/protobuf':'14983d36c56a814ed91b6d652f2b8f895baba1b84eb43b28a0b132c8637cd274',
}
zip_evidence={}
for package in s['packages']:
    if package['name'] != 'gotth-extension-webhook':
        checksum = package['checksums'][0]
        assert checksum['algorithm'] == 'SHA256' and len(checksum['checksumValue']) == 64
        target=f"{package['name']}@{package['versionInfo']}"
        metadata=json.loads(subprocess.check_output(['go','mod','download','-json',target],cwd=root))
        assert checksum['checksumValue']==expected_zip[package['name']]
        assert checksum['checksumValue']==hashlib.sha256(open(metadata['Zip'],'rb').read()).hexdigest()
        assert metadata['Sum'] in package['comment']
        zip_evidence[f"pkg:golang/{target}"]=(checksum['checksumValue'],base64.b64decode(metadata['Sum'][3:]).hex())
with open(provenance, encoding='utf-8') as handle:
    p=json.load(handle)
assert p['_type']=='https://in-toto.io/Statement/v1'
assert p['predicateType']=='https://slsa.dev/provenance/v1'
subjects={item['name']: item['digest']['sha256'] for item in p['subject']}
def digest(path):
    return hashlib.sha256(open(path,'rb').read()).hexdigest()
assert subjects[archive.rsplit('/',1)[-1]]==digest(archive)
assert subjects[sbom.rsplit('/',1)[-1]]==digest(sbom)
assert p['predicate']['buildDefinition']['externalParameters']=={'target':'linux/amd64','version':'1.0.0'}
assert p['predicate']['buildDefinition']['internalParameters']['goVersion'].startswith('go version go1.27.3')
for material in p['predicate']['buildDefinition']['resolvedDependencies'][2:]:
    assert (material['digest']['sha256'],material['digest']['dirHash1']) == zip_evidence[material['uri']]
parameters=p['predicate']['buildDefinition']['internalParameters']
assert parameters['goEnvironment']=={'GOAMD64':'v1','GOENV':'off','GOEXPERIMENT':'','GOFIPS140':'off','GOFLAGS':'','GOPROXY':'off','GOTOOLCHAIN':'local','GOWORK':'off'}
assert 'metadata' not in p['predicate']['runDetails']
PY
grep -Fq 'github.com/gotthboard/gotth-extensions v1.0.0 — MIT' "$extract/$name/THIRD_PARTY_NOTICES.md"
grep -Fq 'github.com/gotthboard/gotth-sdk v1.0.0 — MIT' "$extract/$name/THIRD_PARTY_NOTICES.md"
test "$(find "$extract/$name/licenses" -maxdepth 1 -type f | wc -l)" -ge 8
before=$(sha256sum "$right/$archive" | cut -d' ' -f1)
if SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" "$version" "$right" >/dev/null 2>&1; then
    echo "existing release bundle was replaced" >&2
    exit 1
fi
after=$(sha256sum "$right/$archive" | cut -d' ' -f1)
test "$before" = "$after"
if SOURCE_DATE_EPOCH=$((epoch + 1)) "$root/scripts/build-artifact.sh" "$version" "$right_root/epoch" >/dev/null 2>&1; then
    echo "mismatched source epoch was accepted" >&2
    exit 1
fi
if SOURCE_DATE_EPOCH="$epoch" "$root/scripts/build-artifact.sh" 1.0.0garbage "$right_root/version" >/dev/null 2>&1; then
    echo "invalid release version was accepted" >&2
    exit 1
fi
git clone -q --no-hardlinks "$root" "$clone_root/source"
printf '\nimmutable export regression\n' >> "$clone_root/source/SECURITY.md"
SOURCE_DATE_EPOCH="$epoch" "$clone_root/source/scripts/build-artifact.sh" "$version" "$clone_root/bundle" >/dev/null
mkdir "$clone_root/extract"
tar -xzf "$clone_root/bundle/$archive" -C "$clone_root/extract"
git -C "$clone_root/source" show HEAD:SECURITY.md > "$clone_root/expected-security.md"
cmp "$clone_root/expected-security.md" "$clone_root/extract/$name/SECURITY.md"
python3 - "$clone_root/bundle/$provenance" "$(git -C "$clone_root/source" rev-parse HEAD)" <<'PY'
import json,sys
p=json.load(open(sys.argv[1]))
assert p['predicate']['buildDefinition']['resolvedDependencies'][0]['digest']['gitCommit']==sys.argv[2]
assert p['predicate']['buildDefinition']['internalParameters']['sourceExport']=='git-archive'
assert p['predicate']['buildDefinition']['internalParameters']['invokingCheckoutClean'] is False
PY
git clone -q --no-hardlinks "$root" "$replace_root/source"
mkdir "$replace_root/source/local-protobuf"
printf 'module google.golang.org/protobuf\n\ngo 1.27.3\n' > "$replace_root/source/local-protobuf/go.mod"
printf '\nreplace google.golang.org/protobuf => ./local-protobuf\n' >> "$replace_root/source/go.mod"
git -C "$replace_root/source" add go.mod local-protobuf/go.mod
git -C "$replace_root/source" -c user.name=Replacement -c user.email=replacement@localhost commit -q -m 'test replacement rejection'
if "$replace_root/source/scripts/build-artifact.sh" "$version" "$replace_root/bundle" >"$replace_root/build.log" 2>&1; then
    echo 'committed module replacement was accepted' >&2
    exit 1
fi
grep -Fq 'module replacement is forbidden in release builds' "$replace_root/build.log"
printf '%s\n' 'v1.0.0 immutable-export archive, SPDX SBOM, ZIP and Go h1 checksums, sanitized environment, replacement rejection, and SLSA provenance are reproducible'
