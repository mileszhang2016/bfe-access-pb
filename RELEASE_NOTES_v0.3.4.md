# Release v0.3.4

## What's Changed

- Add AI protocol style field to `RequestLog`:
  - `ai_protocol` (field `717`)

This optional field records the AI request protocol style (e.g., `openai`, `anthropic`) and is intended to support protocol-aware log analysis, routing reconciliation, and provider-specific billing in the AI Gateway.

## Full Changelog

- `bfe_access_pb/bfe_access.proto`: add `ai_protocol` field (717).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 717.
