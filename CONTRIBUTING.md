# Contributing to zeroone

Thanks for your interest in helping out.

## Development setup

You need:

- Go 1.25+ (the toolchain version is pinned in `go.mod`).
- Node.js 20+ and npm (for the web panel).
- Docker 24+ (only required for image / compose work).

```bash
git clone https://github.com/amirrezakm/zeroone.git
cd zeroone
go test ./...
npm --prefix web/app ci
npm --prefix web/app run build
```

`scripts/check.sh` runs the full pre-merge validation pass: regenerates
the local stack config, runs `go test ./...`, renders the Xray config,
and builds the panel if `web/app/node_modules` exists.

## Branches, commits, PRs

- Branch off `main`. Feature branches: `feat/<short-name>`, fixes:
  `fix/<short-name>`.
- Keep commits focused. Squash noise before requesting review.
- Conventional-Commits-style prefixes are welcome but not required.
- Open a pull request against `main`; CI (Go tests, UI build, shellcheck,
  Docker build smoke) must pass before merge.

## Vendoring

`vendor/` is committed so the build works on offline / restricted hosts.
If you add a Go dependency:

```bash
go get example.com/pkg@vX.Y.Z
go mod tidy
go mod vendor
```

Commit `go.mod`, `go.sum`, and the `vendor/` changes together.

## Code style

- `gofmt` clean (`gofmt -l .` is part of CI; the job fails on any diff).
- Prefer the standard library + `log/slog`. No new runtime deps without a
  strong reason — Go's stdlib plus our existing `golang.org/x/net` is
  intentionally small.
- Tests live next to the code: `foo_test.go`.

## Adding a new host-side feature

Anything that shells out to `systemctl`, `nft`, `tc`, `iptables`, or
`journalctl` belongs behind an opt-in flag (`-manage-…`) so the
container build can leave it disabled. See `internal/tunnel/supervisor.go`
for the existing pattern.

## Reporting bugs

Open a GitHub issue with:

- `zeroone version` output.
- `zeroone logs --tail 200`.
- The relevant section of your `stack.json` (with secrets redacted).
- Steps to reproduce.

For security issues, follow [`SECURITY.md`](SECURITY.md) instead.
