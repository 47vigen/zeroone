# Offline Install (Air-Gapped / Iran)

This guide is for operators whose **destination server cannot reach
GHCR, Docker Hub, or GitHub raw**, but who can SFTP into it from
another host. It is the standard workflow for VPS hosts inside Iran
behind upstream filtering.

The flow is a strict three-step:

1. On a **builder host** (anywhere with internet, or with access to
   `mirror-docker.runflare.com`), produce a single tarball containing
   the Docker image, compose file, env example, and an offline
   installer.
2. Transfer that tarball to the destination via **SFTP** (or `scp` /
   `rsync`).
3. On the destination, **load** the image and start the container
   using the bundled installer. No registry pulls happen here.

```
  ┌──────────────┐    1. offline-bundle.sh       ┌──────────────┐
  │  builder     │  ─────────────────────────►   │  tarball     │
  │  (mirror or  │     docker pull + save        │  *.tar.gz    │
  │  internet)   │                               └──────┬───────┘
  └──────────────┘                                      │ 2. sftp
                                                        ▼
                                                ┌──────────────┐
                                                │ destination  │
                                                │ (air-gapped) │
                                                │ docker load  │  3. install-offline.sh
                                                │ compose up   │
                                                └──────────────┘
```

## Prerequisites

**Builder host:**

- Docker (any recent version with `docker compose v2`).
- Outbound access to one of:
  - `mirror-docker.runflare.com` (Iranian PaaS mirror — default), or
  - `ghcr.io` directly (override `ZEROONE_IMAGE_SRC`).
- A checkout of this repository (`git clone` or downloaded tarball).
- `tar` and `sha256sum` (or `shasum`).

**Destination server:**

- Linux x86_64 or aarch64 (Debian, Ubuntu, RHEL family, AlmaLinux,
  Rocky, or Alpine).
- Root access.
- `/var/lib/zeroone` writable (the script creates it).
- **Optional:** access to `mirror-linux.runflare.com` — only needed
  if Docker is not yet installed on the destination and you want the
  script to install it automatically via the Runflare apt mirror.
  Pre-install Docker yourself to skip this entirely.

## Step 1 — Build the bundle on the builder host

```bash
cd zeroone
bash scripts/offline-bundle.sh
```

Output:

```
dist/zeroone-offline-latest-amd64.tar.gz
```

The script prints the path, size, and SHA-256 of the tarball. Save the
SHA-256 — you'll use it to verify the transfer on the destination.

### Tuning the bundle

All knobs are environment variables, set before running the script:

| Variable | Default | Purpose |
| --- | --- | --- |
| `ZEROONE_VERSION` | `latest` | Image tag to pull and bundle (e.g. `v1.0.0`). Pin to a real version in production. |
| `ZEROONE_ARCH` | `amd64` | Target architecture for the destination (`amd64` or `arm64`). Must match `uname -m` on the destination. |
| `ZEROONE_REPO` | `amirrezakm/zeroone` | Source repo path on the registry. Override for forks. |
| `ZEROONE_IMAGE_SRC` | `mirror-docker.runflare.com/$ZEROONE_REPO` | Where to `docker pull` from. Override to `ghcr.io/$ZEROONE_REPO` if your builder has direct GHCR access. |
| `ZEROONE_IMAGE_DST` | `ghcr.io/$ZEROONE_REPO` | The tag the bundled `docker-compose.yml` expects. You should rarely change this. |
| `OUT_DIR` | `./dist` | Where to write the final tarball. |

Example — pin a specific version, target an aarch64 destination, and
pull directly from GHCR:

```bash
ZEROONE_VERSION=v1.0.0 \
ZEROONE_ARCH=arm64 \
ZEROONE_IMAGE_SRC=ghcr.io/amirrezakm/zeroone \
bash scripts/offline-bundle.sh
```

## Step 2 — Transfer to the destination

`sftp` is the path requested by most Iranian operators because it works
through restrictive firewalls that block plain `scp` outbound. Any
file-transfer method that delivers the tarball intact is fine.

```bash
sftp root@destination.example.com
sftp> put dist/zeroone-offline-latest-amd64.tar.gz
sftp> bye
```

On the destination, verify the SHA-256 matches what the bundler
printed:

```bash
sha256sum zeroone-offline-latest-amd64.tar.gz
```

Equivalents: `scp dist/zeroone-offline-*.tar.gz root@host:` or
`rsync -avP dist/zeroone-offline-*.tar.gz root@host:`.

## Step 3 — Install on the destination

```bash
tar -xzf zeroone-offline-latest-amd64.tar.gz
sudo bash install-offline.sh
```

What happens:

1. **Docker check.** If Docker is already installed, the script skips
   ahead. If not, on Debian / Ubuntu it temporarily rewrites
   `/etc/apt/sources.list` to point at
   `http://mirror-linux.runflare.com`, runs `apt-get install
   docker.io docker-compose-v2`, then restores the original
   `sources.list`. On RHEL / Alpine the script stops with an
   instruction to install Docker manually — see "Troubleshooting".
2. **State directories** `/opt/zeroone` and `/var/lib/zeroone` are
   created.
3. **Compose + env** are copied from the bundle to
   `/opt/zeroone/docker-compose.yml` and `/opt/zeroone/.env`.
4. **The image is loaded** with `docker load -i zeroone-image*.tar`.
   No registry traffic.
5. **Admin prompt** — username + password, exactly like the online
   installer. Set `ZEROONE_ADMIN_USERNAME` / `ZEROONE_ADMIN_PASSWORD`
   in the environment to skip the prompt:

   ```bash
   sudo ZEROONE_ADMIN_USERNAME=admin \
        ZEROONE_ADMIN_PASSWORD='choose-a-long-one' \
        bash install-offline.sh
   ```

6. The script self-copies to `/usr/local/bin/zeroone`, brings the
   container up with `docker compose up -d`, waits for `/api/health`
   to respond, and creates the admin account inside the container.

When done, you'll see:

```
 ok zeroone is installed (offline).
   panel:     http://YOUR_HOST_IP:8000/
   listen:    0.0.0.0:8000
   compose:   /opt/zeroone/docker-compose.yml
   state:     /var/lib/zeroone
   cli:       /usr/local/bin/zeroone
```

### Useful flags

```bash
sudo bash install-offline.sh --bundle /path/to/extracted/bundle
sudo bash install-offline.sh --force            # overwrite existing .env
```

### Day-2 operations

Once installed, use the standard `zeroone` CLI:

```bash
zeroone status         # container + panel health
zeroone logs -f        # follow daemon logs
zeroone restart
zeroone cli admin list -config /var/lib/zeroone/stack.json
zeroone backup -o /root/zeroone-backup.tgz
```

`zeroone update` (the online flavor) is **not** available offline.
See the next section.

## Upgrading offline

Repeat the same three steps with a new version:

1. On the builder, rebuild with the new version:
   ```bash
   ZEROONE_VERSION=v1.1.0 bash scripts/offline-bundle.sh
   ```
2. SFTP the new tarball to the destination.
3. On the destination, extract and run the offline `update`
   subcommand against the new bundle directory:
   ```bash
   mkdir -p /root/zeroone-v1.1.0
   tar -xzf zeroone-offline-v1.1.0-amd64.tar.gz -C /root/zeroone-v1.1.0
   sudo zeroone update -b /root/zeroone-v1.1.0
   ```

The update loads the new image, refreshes `docker-compose.yml`,
restarts the container, and prunes the old image. Your `.env`,
`stack.json`, and other state are untouched.

## Troubleshooting

**`docker compose up` says "manifest unknown" or "image not found"**

The bundle was built for a different architecture than the
destination's CPU. Verify with `uname -m` on the destination; rebuild
with `ZEROONE_ARCH=arm64` (or `amd64`) on the builder.

**Docker is missing and `apt-get update` fails**

The Runflare apt mirror requires `apt-transport-https` for HTTPS
sources (the bundle uses plain HTTP, so this is rarely the issue).
If `mirror-linux.runflare.com` itself is unreachable from your
destination, install Docker manually:

```bash
# Debian (example for bookworm; substitute your codename)
cat > /etc/apt/sources.list <<EOF
deb http://mirror-linux.runflare.com/debian bookworm main
deb http://mirror-linux.runflare.com/debian-security bookworm-security main
EOF
apt-get update
apt-get install -y docker.io docker-compose-v2
```

On RHEL/AlmaLinux/Rocky, point `/etc/yum.repos.d/` baseurls at
`https://mirror-linux.runflare.com/almalinux` (or the matching
distribution mirror) and install `docker docker-compose-plugin`. Re-run
`install-offline.sh` once `docker compose version` works.

Set `ZEROONE_SKIP_DOCKER_INSTALL=1` if you want `install-offline.sh`
to refuse to touch `sources.list` under any circumstance.

**`apt-get install` succeeds but `docker compose version` errors**

You have the legacy `docker-compose` (Python, v1) installed instead of
the v2 plugin. Install the plugin: `apt-get install
docker-compose-v2` (Debian/Ubuntu) or `dnf install
docker-compose-plugin` (RHEL family).

**SHA-256 mismatch after SFTP**

The tarball was corrupted in transit. Rebuild and retransfer; if it
keeps happening, try `rsync -avP --checksum` or compress with
`xz` instead.

**`docker exec zeroone zeroone admin add` fails**

The container started but didn't pick up the bundled compose / env, or
the panel hasn't created `/var/lib/zeroone/stack.json` yet. Check
`zeroone logs` and retry the command manually:

```bash
zeroone cli admin add -config /var/lib/zeroone/stack.json \
    -username admin -password 'choose-a-long-one'
```

## Alternative — direct Runflare mirror (no SFTP)

If your destination itself can reach `mirror-docker.runflare.com`, you
can skip SFTP entirely. Configure Docker to use the Runflare registry
mirror, then run the normal online installer.

Create `/etc/docker/daemon.json`:

```json
{
  "registry-mirrors": ["https://mirror-docker.runflare.com"]
}
```

Restart Docker, then run the standard online installer:

```bash
systemctl restart docker
sudo bash -c "$(curl -sSL https://raw.githubusercontent.com/amirrezakm/zeroone/main/scripts/install.sh)" @ install
```

Note that this still requires the destination to reach
`raw.githubusercontent.com` to fetch `docker-compose.yml` and
`.env.example`. If that's blocked too, fall back to the SFTP flow
above.

## See also

- [`docs/INSTALL.md`](INSTALL.md) — the standard online install path.
- [`docs/HOST-INSTALL.md`](HOST-INSTALL.md) — systemd / no-Docker
  install for advanced host integrations.
- [`docs/runflare-edge-deploy.md`](runflare-edge-deploy.md) — using a
  Runflare PaaS edge in front of an origin server.
