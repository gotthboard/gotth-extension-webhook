#!/usr/bin/env python3
"""Generate deterministic release inventory, SBOM, licenses, and provenance."""

import argparse
import base64
import datetime
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess

MODULE_POLICY = {
    "github.com/gotthboard/gotth-extensions": {"version": "v1.0.0", "h1": "h1:lKXOydN/7+Mq9W2juqnWUXNJiirpYR6F3KPXYV55Hho=", "zip_sha256": "f5caba0cb6db41f0b53887316d57450ce7482a78921a12d770d884aaca6ff4e3", "license": "MIT"},
    "github.com/gotthboard/gotth-sdk": {"version": "v1.0.0", "h1": "h1:RmjRs/l+H3KMDmnimZs2OL8ExCchT/opVO5tDdjb4EA=", "zip_sha256": "c975e47838b9882ca1d056fa2882ab220eb3e6aecb8f37ff1feb6776cdc37102", "license": "MIT"},
    "golang.org/x/net": {"version": "v0.60.0", "h1": "h1:79p50tfZlm0J9YfoDsSi639qSXNGVwEzOPLCxM2FsYU=", "zip_sha256": "052e0453db47a96ae9567503e7c95ea925ab52682cb1cec31250541f7160107e", "license": "BSD-3-Clause"},
    "golang.org/x/sys": {"version": "v0.48.0", "h1": "h1:bbX/i/6MgT9BVLM9RT1thmxL04yeTAhbEz4SyadbXoo=", "zip_sha256": "3a47e06431ba1e24d5b2aca3a7ce31fd13faff34df3e5d8e63f105ff9d0c1786", "license": "BSD-3-Clause"},
    "golang.org/x/text": {"version": "v0.42.0", "h1": "h1:JbOZXgfeCPU9gacVtYliJqOhD+zhrEqK4LfdpmlUZqI=", "zip_sha256": "a7b64e003056b6470303f408202098d8f3714a115f23091b8cac85edeb265476", "license": "BSD-3-Clause"},
    "google.golang.org/genproto/googleapis/rpc": {"version": "v0.0.0-20260526163538-3dc84a4a5aaa", "h1": "h1:mZHHdPZl0dbGHCflZgAq/Q468DWVFcU2whhB2KAo8fk=", "zip_sha256": "55f932b106d6f4f8cb76ae869e44f5d6d8400965fa441ee3532d09c96d64aac1", "license": "Apache-2.0"},
    "google.golang.org/grpc": {"version": "v1.83.2", "h1": "h1:EManeRomTObA0BU7I8vXgg/78uE5MJ9M8B39EX2WscU=", "zip_sha256": "84292905a694e17862d15abc123a6f7f7efa9e1e3a77e80120982ff95b2605ea", "license": "Apache-2.0"},
    "google.golang.org/protobuf": {"version": "v1.36.11", "h1": "h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=", "zip_sha256": "14983d36c56a814ed91b6d652f2b8f895baba1b84eb43b28a0b132c8637cd274", "license": "BSD-3-Clause"},
}


def run(args, cwd):
    return subprocess.run(args, cwd=cwd, check=True, stdout=subprocess.PIPE).stdout


def sha256(path):
    digest = hashlib.sha256()
    with open(path, "rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def write_json(path, value):
    encoded = json.dumps(value, sort_keys=True, separators=(",", ":"))
    Path(path).write_text(encoded + "\n", encoding="utf-8")


def go_dirhash_hex(value):
    if not value.startswith("h1:"):
        raise SystemExit(f"unsupported Go module checksum: {value}")
    decoded = base64.b64decode(value[3:], validate=True)
    if len(decoded) != 32:
        raise SystemExit(f"invalid Go module checksum: {value}")
    return decoded.hex()


def module_download(root, module):
    target = f"{module['path']}@{module['version']}"
    metadata = json.loads(run(["go", "mod", "download", "-json", target], root))
    if metadata.get("Error") or metadata.get("Sum") != module["sum"]:
        raise SystemExit(f"Go module checksum mismatch: {target}")
    archive = metadata.get("Zip")
    if not archive or not Path(archive).is_file():
        raise SystemExit(f"cached Go module ZIP unavailable: {target}")
    archive_digest = sha256(archive)
    expected = MODULE_POLICY[module["path"]]
    if module["version"] != expected["version"] or module["sum"] != expected["h1"] or archive_digest != expected["zip_sha256"]:
        raise SystemExit(f"reviewed module identity mismatch: {target}")
    return archive, archive_digest


def spdx_id(value):
    return "SPDXRef-" + re.sub(r"[^A-Za-z0-9.-]", "-", value)


def linked_modules(binary, root):
    output = run(["go", "version", "-m", binary], root).decode("utf-8")
    modules = []
    for line in output.splitlines():
        fields = line.split("\t")
        if len(fields) >= 2 and fields[1] == "=>":
            raise SystemExit("replacement module linked into release binary")
        if len(fields) >= 5 and fields[1] == "dep":
            modules.append({"path": fields[2], "version": fields[3], "sum": fields[4]})
    modules.sort(key=lambda item: item["path"])
    actual = {item["path"] for item in modules}
    expected = set(MODULE_POLICY)
    if actual != expected:
        missing = sorted(expected - actual)
        extra = sorted(actual - expected)
        raise SystemExit(f"linked module allowlist mismatch: missing={missing} extra={extra}")
    for module in modules:
        module["zip"], module["zip_sha256"] = module_download(root, module)
    return modules


def source_manifest(root):
    root_path = Path(root)
    files = sorted(path for path in root_path.rglob("*") if path.is_file())
    lines = []
    for path in files:
        relative = path.relative_to(root_path).as_posix()
        lines.append(f"{sha256(path)}  {relative}\n")
    manifest = "".join(lines)
    return manifest, hashlib.sha256(manifest.encode("utf-8")).hexdigest()


def module_directory(root, module):
    target = f"{module['path']}@{module['version']}"
    return Path(run(["go", "list", "-m", "-f={{.Dir}}", target], root).decode().strip())


def inventory(args):
    root = os.path.realpath(args.root)
    output = Path(args.output)
    output.mkdir(parents=True, exist_ok=True)
    licenses = output / "licenses"
    licenses.mkdir(mode=0o755)
    modules = linked_modules(args.binary, root)
    manifest, source_digest = source_manifest(root)
    (output / "SOURCE_MANIFEST.sha256").write_text(manifest, encoding="utf-8")

    notice_lines = [
        "# Third-party notices\n\n",
        "This binary links the following Go modules. License and notice files are ",
        "included verbatim under licenses/. The SPDX identifier is an admission ",
        "policy classification; the copied text remains authoritative.\n\n",
    ]
    packages = [{
        "SPDXID": "SPDXRef-Package-gotth-extension-webhook",
        "name": "gotth-extension-webhook",
        "versionInfo": args.version,
        "downloadLocation": "NOASSERTION",
        "filesAnalyzed": False,
        "licenseConcluded": "MIT",
        "licenseDeclared": "MIT",
        "copyrightText": "Copyright (c) 2026 GOTTH Board contributors",
        "externalRefs": [{
            "referenceCategory": "PACKAGE-MANAGER",
            "referenceType": "purl",
            "referenceLocator": f"pkg:golang/github.com/gotthboard/gotth-extension-webhook@{args.version}",
        }],
        "comment": f"Tracked-source manifest SHA-256: {source_digest}",
    }]
    relationships = []
    for module in modules:
        identifier = spdx_id("Package-" + module["path"])
        license_id = MODULE_POLICY[module["path"]]["license"]
        directory = module_directory(root, module)
        candidates = sorted(
            path for path in directory.iterdir()
            if path.is_file() and path.name.upper().startswith(("LICENSE", "COPYING", "NOTICE"))
        )
        if not any(path.name.upper().startswith(("LICENSE", "COPYING")) for path in candidates):
            raise SystemExit(f"no license file for linked module {module['path']} {module['version']}")
        copied = []
        prefix = re.sub(r"[^A-Za-z0-9_.-]", "_", module["path"] + "@" + module["version"])
        for source in candidates:
            target = licenses / f"{prefix}.{source.name}"
            shutil.copyfile(source, target)
            copied.append("licenses/" + target.name)
        notice_lines.append(
            f"- {module['path']} {module['version']} — {license_id}; "
            + ", ".join(copied) + "\n"
        )
        packages.append({
            "SPDXID": identifier,
            "name": module["path"],
            "versionInfo": module["version"],
            "downloadLocation": f"https://proxy.golang.org/{module['path']}/@v/{module['version']}.zip",
            "filesAnalyzed": False,
            "checksums": [{"algorithm": "SHA256", "checksumValue": module["zip_sha256"]}],
            "licenseConcluded": license_id,
            "licenseDeclared": license_id,
            "copyrightText": "NOASSERTION",
            "externalRefs": [{
                "referenceCategory": "PACKAGE-MANAGER",
                "referenceType": "purl",
                "referenceLocator": f"pkg:golang/{module['path']}@{module['version']}",
            }],
            "comment": f"Go module checksum {module['sum']}; license files: {', '.join(copied)}",
        })
        relationships.append({
            "spdxElementId": "SPDXRef-Package-gotth-extension-webhook",
            "relationshipType": "DEPENDS_ON",
            "relatedSpdxElement": identifier,
        })
    (output / "THIRD_PARTY_NOTICES.md").write_text("".join(notice_lines), encoding="utf-8")
    created = datetime.datetime.fromtimestamp(args.epoch, datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    sbom = {
        "spdxVersion": "SPDX-2.3",
        "dataLicense": "CC0-1.0",
        "SPDXID": "SPDXRef-DOCUMENT",
        "name": f"gotth-extension-webhook-{args.version}-linux-amd64",
        "documentNamespace": f"https://gotthboard.org/spdx/gotth-extension-webhook/{args.version}/{source_digest}",
        "creationInfo": {"created": created, "creators": ["Tool: gotth-extension-webhook-release-metadata-v1"]},
        "documentDescribes": ["SPDXRef-Package-gotth-extension-webhook"],
        "packages": packages,
        "relationships": relationships,
    }
    write_json(output / "SBOM.spdx.json", sbom)


def provenance(args):
    root = os.path.realpath(args.root)
    modules = linked_modules(args.binary, root)
    source_digest = sha256(args.source_manifest)
    archive_digest = sha256(args.archive)
    sbom_digest = sha256(args.sbom)
    go_version = run(["go", "version"], root).decode("utf-8").strip()
    materials = [{
        "uri": f"git+https://git.dannyhunn.com/gotthboard/gotth-extension-webhook.git@{args.commit}",
        "digest": {"gitCommit": args.commit, "gitTree": args.tree},
    }, {
        "uri": "https://gotthboard.org/source-manifest/gotth-extension-webhook",
        "digest": {"sha256": source_digest},
    }]
    for module in modules:
        materials.append({
            "uri": f"pkg:golang/{module['path']}@{module['version']}",
            "digest": {"sha256": module["zip_sha256"], "dirHash1": go_dirhash_hex(module["sum"])},
        })
    statement = {
        "_type": "https://in-toto.io/Statement/v1",
        "subject": [
            {"name": os.path.basename(args.archive), "digest": {"sha256": archive_digest}},
            {"name": os.path.basename(args.sbom), "digest": {"sha256": sbom_digest}},
        ],
        "predicateType": "https://slsa.dev/provenance/v1",
        "predicate": {
            "buildDefinition": {
                "buildType": "https://git.dannyhunn.com/gotthboard/gotth-extension-webhook/src/branch/main/docs/distribution.md#reproducible-release-build",
                "externalParameters": {"version": args.version, "target": "linux/amd64"},
                "internalParameters": {
                    "cgoEnabled": False,
                    "goVersion": go_version,
                    "goBuildFlags": [
                        "-trimpath",
                        "-buildvcs=false",
                        f"-ldflags=-s -w -buildid= -X github.com/gotthboard/gotth-extension-webhook/internal/webhook.BuildVersion={args.version}",
                    ],
                    "sourceExport": "git-archive",
                    "invokingCheckoutClean": args.clean,
                    "goEnvironment": {"GOENV": "off", "GOWORK": "off", "GOFLAGS": "", "GOTOOLCHAIN": "local", "GOEXPERIMENT": "", "GOFIPS140": "off", "GOAMD64": "v1", "GOPROXY": "off"},
                },
                "resolvedDependencies": materials,
            },
            "runDetails": {
                "builder": {"id": f"git+https://git.dannyhunn.com/gotthboard/gotth-extension-webhook.git@{args.commit}#scripts/build-artifact.sh"},
            },
        },
    }
    write_json(args.output, statement)


def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(dest="command", required=True)
    inv = subparsers.add_parser("inventory")
    inv.add_argument("--root", required=True)
    inv.add_argument("--version", required=True)
    inv.add_argument("--binary", required=True)
    inv.add_argument("--output", required=True)
    inv.add_argument("--epoch", required=True, type=int)
    inv.set_defaults(function=inventory)
    prov = subparsers.add_parser("provenance")
    prov.add_argument("--root", required=True)
    prov.add_argument("--version", required=True)
    prov.add_argument("--binary", required=True)
    prov.add_argument("--archive", required=True)
    prov.add_argument("--sbom", required=True)
    prov.add_argument("--source-manifest", required=True)
    prov.add_argument("--commit", required=True)
    prov.add_argument("--tree", required=True)
    prov.add_argument("--clean", action="store_true")
    prov.add_argument("--output", required=True)
    prov.set_defaults(function=provenance)
    args = parser.parse_args()
    args.function(args)


if __name__ == "__main__":
    main()
