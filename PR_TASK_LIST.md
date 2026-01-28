# PR Completion Task List

This task list captures the remaining work to satisfy the PR checklist and acceptance criteria for **Public Video Download + Progress UI (Terminal Tool)**.

## A) Scope, Compliance, and Non-Goals (Hard Gates)

- [x] Ensure only **publicly accessible** URLs are supported (no auth/paywall/DRM content).
- [x] Verify there is **no bypass logic** (DRM circumvention, credential harvesting, cookie reuse, or browser automation to access restricted streams).
- [x] Add explicit **restricted content detection** with actionable errors and non-zero exit codes for login/paywall/DRM.

## B) Source & Format Support

- [x] Accept any valid public video URL (no hardcoded allowlist).
- [x] Support direct file downloads for `.mp4`, `.webm`, `.mov` (or equivalent).
- [x] Support **unencrypted** HLS (`.m3u8`) and DASH (`.mpd`).
- [x] Detect and reject encrypted/DRM manifests (HLS AES-128 key URIs, DASH Widevine/PlayReady/CENC).
- [x] Keep source handling modular/extensible (pluggable extractors/parsers).

## C) URL Analysis & Validation

- [x] Validate URL format with explicit errors on invalid input.
- [x] Detect downloadability before downloading (public access + format).
- [x] Implement `--list-formats` for multi-variant sources.

## D) Download Behavior & Output Correctness

- [x] Default to best available quality.
- [x] Support `--quality` (resolution/bitrate) and `--format` (container/codec) selection.
- [x] Support `--output` templates/paths.
- [x] Download + assemble segmented streams in order.
- [x] Support resume without corruption.
- [x] Validate playable, correctly muxed output.

## E) CLI Interface Requirements

- [x] Terminal-only operation with `url` as required input.
- [x] Optional flags: `--quality`, `--format`, `--output`, `--list-formats`, `--audio-only` (or explicitly omit), `--json`.
- [x] `--json` emits **only JSON** to stdout (no noise/progress).

## F) Error Handling & Messaging

- [x] Categorize errors: invalid URL, unsupported format, restricted access/DRM, network failures, filesystem errors.
- [x] Make messages actionable with consistent non-zero exit codes.

## G) Performance & Robustness

- [x] Stream I/O (no excessive memory use).
- [x] Parallel segment downloads where applicable.
- [x] Throttle progress updates (no busy looping).
- [x] Predictable concurrency scaling.

## H) Security Requirements

- [x] Sanitize output paths/filenames (prevent traversal/injection).
- [x] Never execute downloaded content.
- [x] Avoid storing credentials/cookies.
- [x] Use safe timeouts and retry strategies for requests.

## I) Progress UI Integration (Hard Requirement)

- [x] Support user-defined layouts (CLI/config) with fields (label, %, rate, ETA, bytes).
- [x] Multiple progress bars: stable ordering, no flicker/corruption.
- [x] Handle terminal resize with reflow and preserved alignment.
- [x] Interleaved logging: logs above bars, render resumes cleanly.
- [x] Controlled refresh rate; low overhead.
- [x] Fallback for non-TTY/ANSI unsupported (plain text).
- [x] Completed bars persist until summary (or consistent cleanup).
- [x] Renderer decoupled via events; renderer must not block I/O.
- [x] Compatible with parallel downloads and resume.

## J) Metadata Collection & Playlist Metadata (Hard Requirement)

- [x] Collect per-item metadata (title, artists, album, track/disc, release date/year, duration, thumbnail, source URL/ID, extractor name/version).
- [x] Playlist metadata: title/id/url, stable ordering, `playlist.json` manifest with positions.
- [x] Implement tiered metadata strategy (platform structured data, oEmbed/OG tags, manifest hints, user overrides).
- [x] Provide `--meta` overrides with documented precedence.
- [x] Fail gracefully when metadata missing; include structured warnings and safe fallback naming.
- [x] Emit sidecar JSON for each item; embed tags into audio when supported, otherwise log skip.

## K) Testing Requirements

- [x] Unit tests: URL parsing/validation, format selection, DRM detection, metadata parsing, playlist ordering.
- [x] Integration tests with public non-restricted sources.
- [x] Progress UI tests: multiple bars, resize handling, logging during progress, non-TTY behavior.
- [x] Validate stable output (not just exit codes).

## L) Documentation Requirements

- [x] Document supported formats and limitations.
- [x] Document all CLI flags with examples.
- [x] Include legal/copyright notice.
- [x] State non-goals explicitly (no login/paywall/DRM support).
- [x] Document metadata behavior (fields, sources, overrides, sidecar JSON schema).

## Acceptance Criteria Checklist

- [x] AC-1: Direct file download (public URL).
- [x] AC-2: Unencrypted HLS download.
- [x] AC-3: Unencrypted DASH download.
- [x] AC-4: Restricted content detection (login/paywall).
- [x] AC-5: DRM/encrypted stream refusal.
- [x] AC-6: `--list-formats` output.
- [x] AC-7: Quality selection.
- [x] AC-8: Resume support.
- [x] AC-9: Multiple concurrent downloads.
- [x] AC-10: Multiple progress bars render correctly.
- [x] AC-11: User-defined progress layout.
- [x] AC-12: Terminal resize handling.
- [x] AC-13: Interleaved logging with progress.
- [x] AC-14: Non-TTY behavior.
- [x] AC-15: JSON output mode cleanliness.
- [x] AC-16: Path safety.
- [x] AC-17: Tests and docs exist.
- [x] AC-18: Metadata collected when available.
- [x] AC-19: Playlist ordering and manifest.
- [x] AC-20: Metadata missing -> graceful degradation.
- [x] AC-21: User metadata overrides.
- [x] AC-22: Metadata embedding behavior.
