# AGENTS.md

## Cursor Cloud specific instructions

### Services

| Service | Command | Port |
|---------|---------|------|
| Go daemon | `go run ./cmd/zeroone -config config/stack.example.json` | 8091 |
| React panel (dev) | `npm --prefix web/app run dev` | 5173 |

The Vite dev server proxies `/api` requests to `127.0.0.1:8091`.

### Running the daemon locally

- The daemon listens on `127.0.0.1:8091` by default (set in `stack.example.json`).
- On first run with no admins configured, the API is in "bootstrap" mode — POST `/api/admins` with `{"username":"…","password":"…"}` to create the first admin (no auth required).
- The `/var/lib/zeroone` mkdir warning is expected when running as non-root; does not affect functionality.
- Do **not** pass `-allow-apply` in local dev — it attempts to mutate live Xray/systemd state.

### Lint / test / build

See `CLAUDE.md` and `CONTRIBUTING.md` for full details. Quick reference:

- **Go tests:** `go test ./...` (or `go test -race -count=1 ./...` for CI parity)
- **Go vet:** `go vet ./...`
- **Go format check:** `gofmt -l . | grep -v '^vendor/'`
- **UI lint:** `npm --prefix web/app run lint`
- **UI format check:** `npm --prefix web/app run format:check`
- **UI build:** `npm --prefix web/app run build`
- **Security:** `govulncheck ./...` (install: `go install golang.org/x/vuln/cmd/govulncheck@latest`)

### Vendoring

`vendor/` is committed. After any Go dependency change:
```
go mod tidy && go mod vendor
```
Commit `go.mod`, `go.sum`, and `vendor/` together.

### Gotchas

- The daemon mutates `config/stack.example.json` at runtime (adds admin hashes, session secrets, subscription tokens). Use `git checkout -- config/stack.example.json` before committing to avoid leaking test credentials.
- `golangci-lint` in CI uses v2.12.2 — match this version if running locally.
- The CI workflow `security.yml` runs `govulncheck`, Trivy filesystem scan, and Trivy container image scan. The `go` directive in `go.mod` must match the latest Go patch release to avoid stdlib vulnerability reports.
