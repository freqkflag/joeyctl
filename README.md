# joeyctl

Homelab helper CLI for Traefik, Cloudflared, and other cluster chores.

## Installation

### From source

```bash
git clone https://github.com/freqkflag/joeyctl.git
cd joeyctl
go build -o joeyctl .
sudo mv joeyctl /usr/local/bin/joeyctl
```

### From releases

1. Download the latest `joeyctl-linux-<arch>` asset from the [releases page](https://github.com/freqkflag/joeyctl/releases).
2. Copy it to `/usr/local/bin/joeyctl`.
3. `chmod +x /usr/local/bin/joeyctl`.
4. Run `joeyctl --help` to verify.

### Via install script

```bash
curl -fsSL https://raw.githubusercontent.com/freqkflag/joeyctl/main/install.sh | sudo bash
```

## Usage

### Traefik catchall helpers

Create/update a wildcard catchall that routes `*.DOMAIN` and `DOMAIN`:

```bash
sudo joeyctl traefik catchall wildcard \
  --backend http://192.168.12.24:80 \
  --domain cultofjoey.com \
  --priority 1
```

Create/update an `anyhost` catchall that matches any Host header:

```bash
sudo joeyctl traefik catchall anyhost \
  --backend http://192.168.12.24:80
```

Remove previously generated catchall definitions:

```bash
sudo joeyctl traefik catchall remove
```

Pass `--dry-run` to preview file changes and `systemctl` invocations without modifying the host.

### Cloudflared helpers

View the available Cloudflared automation commands (work in progress):

```bash
joeyctl cloudflared --help
```

## Releasing

Tagging a commit that matches `v*` (e.g. `v0.1.0`) triggers the GitHub Actions workflow that builds Linux amd64/arm64 binaries and attaches them to the release.
