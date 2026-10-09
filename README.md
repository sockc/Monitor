# Monitor

A lightweight, **Docker-free**, self-hosted Linux server monitor. V0.1 is an initial working MVP.

**Stack:** Go standard library, systemd, embedded dashboard. No Docker, Node.js, Redis, or database service required on hosts.

## Features

- Linux CPU, memory, root filesystem, total network traffic, uptime
- Multiple named agents; outbound HTTPS reporting every 5 seconds
- Private browser dashboard with admin token login and 5-second auto-refresh
- Server-side agent bearer-token authentication
- Atomic JSON snapshot persistence (every 30 seconds); **history charts and SQLite are planned for future releases**
- GitHub Actions builds for linux/amd64, linux/arm64; server and agent use the same binary with different mode flags

## Run server

Download the matching binary from a GitHub Actions build / release (when available), or build on a development machine:

```sh
go build -trimpath -ldflags="-s -w" -o monitor ./cmd/monitor
MONITOR_AGENT_TOKEN="$(openssl rand -hex 32)"  # save this token for the agents
MONITOR_ADMIN_TOKEN="$(openssl rand -hex 32)"  # save separately
sudo mkdir -p /var/lib/monitor
sudo env MONITOR_AGENT_TOKEN="$MONITOR_AGENT_TOKEN" MONITOR_ADMIN_TOKEN="$MONITOR_ADMIN_TOKEN" ./monitor -mode server -listen 127.0.0.1:8090
```

Place an HTTPS reverse proxy (Nginx/Caddy) in front of `127.0.0.1:8090`. Do not expose the HTTP port to the public internet. Open the proxy URL in your browser and log in using the admin token.

## Run agent

```sh
MONITOR_AGENT_TOKEN="<the server's agent token>" ./monitor -mode agent -server https://monitor.example.com -name server-1
```

**Security:** Use a dedicated long random ingestion token and a different admin token. The MVP uses one shared ingestion token across agents: rotating it affects all agents. Cookies contain the admin token and must be protected by HTTPS. Use a private deployment first; per-node credentials, session hardening, encrypted history, and monitoring alerts are future work.

## systemd

See `deploy/systemd` for example units and `scripts/install.sh` for installation guidance. Review local secrets before enabling services.

## Roadmap

- V0.1: Agent reporting, live dashboard, simple persistence, release automation
- V0.2: SQLite, history retention and charts, per-agent token provisioning
- V0.3: HTTP/TCP checks and alerting
- V0.4: safe signed upgrades, backups, detailed access control

## License

MIT.


## V0.2: history and node-scoped tokens

Server now stores samples in SQLite at `/var/lib/monitor/monitor.db` (WAL mode, rolling 30-day history). Use the dashboard's history range and metric selectors for charts. Migration preserves the old snapshot if present, while the new SQLite database becomes the main persistent metric history.

To create or rotate a **per-node** token, log into the admin dashboard first, then send:

```sh
curl -b cookies.txt -X POST https://monitor.example.com/api/v1/tokens -H 'Content-Type: application/json' -d '{"name":"server-1"}'
```

The endpoint requires the authenticated `monitor_session` cookie; tokens are only returned once. The legacy shared `MONITOR_AGENT_TOKEN` remains a fallback **only for nodes without an issued token**. After every node has its own token, disable the shared fallback in a future breaking change. Each token is scoped to its exact agent name; give the agent `-name server-1` when using its token. Ensure transport uses HTTPS.

## V0.2 installation from tagged GitHub Release

After [the V0.2 release](https://github.com/sockc/Monitor/releases/tag/v0.2.0) is published, on Linux AMD64/ARM64:

```sh
curl -fsSL https://raw.githubusercontent.com/sockc/Monitor/main/scripts/install-release.sh -o install-monitor.sh
sudo bash install-monitor.sh server
```

For an agent, set `MONITOR_SERVER` and `MONITOR_AGENT_TOKEN` in the environment before installation, or configure `/etc/monitor/agent.env` before restarting the service. The server installs to localhost-only port 8090; configure HTTPS reverse proxy separately.

Server credentials: `sudo cat /etc/monitor/server.env`. Treat them as secrets. Use `journalctl -u monitor-server -f` / `journalctl -u monitor-agent -f` for diagnostics.
