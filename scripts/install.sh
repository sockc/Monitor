#!/usr/bin/env bash
set -euo pipefail
if [[ "${EUID}" -ne 0 ]]; then echo "Run as root" >&2; exit 1; fi
MODE="${1:-}"
if [[ "$MODE" != server && "$MODE" != agent ]]; then echo "Usage: sudo bash scripts/install.sh server|agent (MONITOR_BINARY=/path/to/monitor)" >&2; exit 1; fi
BIN="${MONITOR_BINARY:-./monitor}"
test -f "$BIN" || { echo "Binary not found: $BIN" >&2; exit 1; }
install -m 0755 "$BIN" /usr/local/bin/monitor
install -d -m 0700 /etc/monitor
if [[ "$MODE" == server ]]; then
 id -u monitor &>/dev/null || useradd --system --home-dir /var/lib/monitor --shell /usr/sbin/nologin monitor
 install -d -o monitor -g monitor -m 0700 /var/lib/monitor
 ENVFILE=/etc/monitor/server.env
 if [[ ! -e "$ENVFILE" ]]; then
  umask 077
  printf 'MONITOR_AGENT_TOKEN=%s\nMONITOR_ADMIN_TOKEN=%s\n' "$(openssl rand -hex 32)" "$(openssl rand -hex 32)" >"$ENVFILE"
 fi
else
 id -u monitor-agent &>/dev/null || useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin monitor-agent
 ENVFILE=/etc/monitor/agent.env
 if [[ ! -e "$ENVFILE" ]]; then
  umask 077
  cat >"$ENVFILE" <<'EOF'
MONITOR_SERVER=https://monitor.example.com
MONITOR_AGENT_TOKEN=REPLACE_WITH_SERVER_AGENT_TOKEN
EOF
  echo "EDIT $ENVFILE before starting the agent."
  exit 0
 fi
fi
chmod 600 "$ENVFILE"
cat >"/etc/systemd/system/monitor-$MODE.service" <<EOF
[Unit]
Description=Monitor $MODE
After=network-online.target
Wants=network-online.target
[Service]
User=$( [[ "$MODE" == server ]] && echo monitor || echo monitor-agent )
EnvironmentFile=$ENVFILE
ExecStart=/usr/local/bin/monitor -mode $MODE $( [[ "$MODE" == server ]] && echo '-listen 127.0.0.1:8090' )
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
$( [[ "$MODE" == server ]] && echo 'StateDirectory=monitor' )
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable --now "monitor-$MODE"
echo "Installed monitor-$MODE. Use: systemctl status monitor-$MODE"
