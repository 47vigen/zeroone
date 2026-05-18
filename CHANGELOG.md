# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Public release preparation: AGPL-3.0 license, CONTRIBUTING / SECURITY /
  CODE_OF_CONDUCT, issue and PR templates.
- Sanitized example configs (`config/stack.example.json`,
  `config/examples/config.example.json`) using RFC-5737 TEST-NET addresses
  and `example.com` hostnames.
- `deploy/skeleton/xray.service` — reference systemd unit for host installs.

### Changed

- Go module path: `github.com/sakhtar/xray-stack-zeroone` →
  `github.com/amirrezakm/zeroone`.
- Repo-root `Dockerfile` (edge-relay for Runflare) moved to
  `cmd/edge-relay/Dockerfile.runflare`.

### Removed

- `rootfs/` — live-host snapshot. Not appropriate for a public repo.
- `scripts/sync-from-server.sh`, `scripts/import-live-stack.py`,
  `scripts/install-local-layout.sh` — internal tooling.
- `docs/migration-cutover.md`, `docs/local-file-inventory.txt`,
  `docs/rewrite-plan.md` — internal history.

## [0.1.0] — TBD

First public release.

[Unreleased]: https://github.com/amirrezakm/zeroone/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/amirrezakm/zeroone/releases/tag/v0.1.0
