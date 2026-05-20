# Security policy

## Supported versions

`zeroone` is currently pre-1.0. Security fixes land on `main` and are
included in the next tagged release. Pinning to an older tag means you
do not receive security updates — pull `latest` (or rebase your fork) to
stay current.

## Reporting a vulnerability

**Do not open a public GitHub issue for security problems.**

Use one of:

- GitHub private vulnerability reporting:
  <https://github.com/amirrezakm/zeroone/security/advisories/new>
- Direct email to the maintainer listed at the top of the repository
  profile.

Please include:

- A description of the issue and its impact.
- Steps to reproduce or a proof-of-concept.
- Affected version(s) (`zeroone version`).
- Whether you've already disclosed this anywhere else.

We aim to acknowledge within 72 hours and ship a fix within 14 days for
high-severity issues. Coordinated disclosure is welcome.

## Hardening recommendations

The admin API exposes full control over the Xray config, user
credentials, and traffic counters. Treat it like an SSH endpoint:

- **Never expose the admin port directly to the public internet.** Put
  it behind a reverse proxy (Nginx, Caddy) with TLS, basic auth, or an
  IP allowlist. Bind to `127.0.0.1` and tunnel via SSH if possible.
- **Use strong admin passwords.** Hashes are stored as PBKDF2-SHA256 with
  per-password salts, but a weak password is still a weak password.
- **Rotate the session secret** by deleting `panel.session_secret` from
  `stack.json` and restarting — invalidates all existing sessions.
- **Keep the host updated.** `unattended-upgrades` on Debian/Ubuntu,
  `dnf-automatic` on RHEL-family, etc.
- **Audit log** (`/var/lib/zeroone/audit.log`) records every mutating
  action. Ship it off-host.
- **Backup**: `zeroone backup -o /root/zeroone-$(date +%F).tgz`
  contains `stack.json` (with password hashes) — store backups encrypted.

See [`docs/SECURITY.md`](docs/SECURITY.md) for the operator-facing
hardening guide.
