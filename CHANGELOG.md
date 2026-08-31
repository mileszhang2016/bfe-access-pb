# Changelog

All notable changes to `bfe-access-pb` will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.5]

### Added

- Add image input token and video metering fields to `RequestLog`:
  - `ai_image_input_tokens` (field `786`)
  - `ai_video_count` (field `787`)

These optional fields support AI image-aware billing where image input tokens are priced separately from text tokens, and video generation billing where cost is calculated per generated video.

### Changed

- `bfe_access_pb/bfe_access.proto`: add image input token and video count fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 786 and 787.

## [v0.3.4]

### Added

- Add AI protocol style field to `RequestLog`:
  - `ai_protocol` (field `717`)

This optional field records the AI request protocol style (e.g., `openai`, `anthropic`) and is intended to support protocol-aware log analysis, routing reconciliation, and provider-specific billing in the AI Gateway.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_protocol` field (717).
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document field 717.

## [v0.3.3]

### Added

- Add AI request mode field to `RequestLog`:
  - `ai_mode` (field `716`)

- Add image metering field to `RequestLog`:
  - `ai_image_count` (field `785`)

These optional fields support AI image generation billing where cost is calculated per generated image, and enable mode-based log analysis and cost reconciliation.

### Changed

- `bfe_access_pb/bfe_access.proto`: add `ai_mode` and `ai_image_count` fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 716 and 785.
- `b2log/test_data/`: add missing test fixtures and `gen_test_data.go`.

## [v0.3.2]

### Added

- Add audio token metering fields to `RequestLog`:
  - `ai_audio_input_tokens` (field `783`)
  - `ai_audio_output_tokens` (field `784`)

These optional fields support AI audio billing where audio tokens are priced separately from text tokens.

### Changed

- `bfe_access_pb/bfe_access.proto`: add audio token fields.
- `bfe_access_pb/bfe_access.pb.go`: regenerate.
- `docs/protobuf.md`: document fields 783 and 784.
