# AGENTS.md — BFE Access PB

This file guides AI coding agents working on the `bfe-access-pb/` repository.

## Project overview

`bfe-access-pb` is a standalone Go module that hosts:

- The Protocol Buffers schema for BFE access logs (`bfe_access_pb/bfe_access.proto`).
- The Go code generated from that schema (`bfe_access_pb/bfe_access.pb.go`).
- The `b2log` binary log wrapper library for reading and writing serialized `BfeLog` records.

Other repositories (notably `bfe/`) import this module to consume the protobuf definitions and the binary log helpers.

## Directory structure

| Directory | Responsibility |
|-----------|----------------|
| `bfe_access_pb/` | Protobuf schema and generated Go code for access logs. |
| `b2log/` | Binary log record reader/writer library. |
| `docs/` | Human-readable documentation: `protobuf.md` and `README.md`. |
| `build.sh` | Script to regenerate `bfe_access.pb.go` from the `.proto` file. |

## Build conventions

- **Go version**: 1.22 (`go.mod`).
- **Module**: `github.com/bfenetworks/bfe-access-pb`.
- **Protobuf generation**: run `sh build.sh` after any change to `bfe_access_pb/bfe_access.proto`.
  - Uses `protoc` (defaults to `/opt/protoc`, override with `PROTOC` env var).
  - Installs `protoc-gen-go@v1.35.0` and regenerates `bfe_access_pb/bfe_access.pb.go` with `paths=source_relative`.
- **Tests**: `go test ./...` from the repository root.

## Common modification patterns

### Add or modify access-log fields

1. Edit `bfe_access_pb/bfe_access.proto`.
   - Pick the next available field number in the appropriate semantic region.
   - AI Observability fields live in the 701–900 range (see `docs/protobuf.md` for sub-ranges).
   - Mark fields as `optional` unless there is a strong reason to require them.
2. Update `docs/protobuf.md` to document the new fields, including field number, type, and semantics.
3. Run `sh build.sh` to regenerate `bfe_access_pb/bfe_access.pb.go`.
4. Verify the generated Go code contains the expected fields/getters.
5. Run `go test ./...` to ensure nothing is broken.
6. Bump the module version tag when the change is ready for downstream consumption (e.g., `v0.3.1`).

### Modify b2log behavior

1. Edit sources under `b2log/`.
2. Add or update unit tests next to the code under test (`*_test.go`).
3. Run `go test ./b2log/...`.

## Agent guidelines

- **Never hand-edit generated files** (`bfe_access.pb.go`). Always regenerate via `build.sh`.
- **Keep docs in sync**: any `.proto` field change must be reflected in `docs/protobuf.md`.
- **License headers**: all new source files need the Apache 2.0 header. Existing files in this repo already carry the header; preserve it.
- **Field number stability**: once a field number is used in a released tag, do not reuse or repurpose it. If a field is no longer needed, `reserve` its number in the proto instead.
- **Version bumps**: this repository is consumed by `bfe/` as a tagged module. After a protobuf or API change, plan for a new tag; during local development in a monorepo, downstream can temporarily use a `replace` directive.

## Useful references

- `docs/protobuf.md` — full field listing and numbering conventions.
- `docs/README.md` — b2log binary format and usage examples.
- `build.sh` — protobuf generation script.
