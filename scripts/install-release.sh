#!/usr/bin/env bash
set -euo pipefail
if [[ "${EUID}" != 0 ]]; then echo "Please run as root" >&2;exit 1;fi
MODE="${1:-}";if [[ "$MODE" != "server" && "$MODE" != "agent" ]];then echo "Usage: sudo bash install-release.sh server|agent" >&2;exit 1;fi
ARCH="$(uname -m)"
case "$ARCH" in x86_64) ARCH=amd64;; aarch64|arm64) ARCH=arm64;; *) echo "Unsupported CPU $ARCH";exit 1;;esac
VERSION="${MONITOR_VERSION:-v0.7.1}"
BASE="https://github.com/sockc/Monitor/releases/download/${VERSION}"
TMP="$(mktemp -d)";trap 'rm -rf "$TMP"' EXIT
curl -fLSs --retry 3 "$BASE/monitor-linux-$ARCH" -o "$TMP/monitor"
curl -fLSs --retry 3 "$BASE/SHA256SUMS" -o "$TMP/SHA256SUMS"
(cd "$TMP";grep "monitor-linux-$ARCH$" SHA256SUMS | sed -E "s@[^[:space:]]*monitor-linux-$ARCH@monitor@" | sha256sum -c -)
install -m 0755 "$TMP/monitor" /usr/local/bin/monitor
install -d -m 0700 /etc/monitor
if [[ "$MODE" == server ]];then
 id monitor >/dev/null 2>&1 || useradd --system --home-dir /var/lib/monitor --shell /usr/sbin/nologin monitor
 install -d -m 0700 -o monitor -g monitor /var/lib/monitor
 ENVFILE=/etc/monitor/server.env
 if ! grep -q "^MONITOR_LISTEN=" "$ENVFILE" 2>/dev/null;then printf "MONITOR_LISTEN=127.0.0.1:%s\\n" "${MONITOR_PORT:-8090}" >> "$ENVFILE";fi
 if [[ ! -f "$ENVFILE" ]];then
  umask 077
  printf 'MONITOR_AGENT_TOKEN=%s\nMONITOR_ADMIN_TOKEN=%s\n' "$(openssl rand -hex 32)" "$(openssl rand -hex 32)" > "$ENVFILE"
 fi
 USERNAME=monitor
 CMD='/usr/local/bin/monitor -mode server'
else
 id monitor-agent >/dev/null 2>&1 || useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin monitor-agent
 ENVFILE=/etc/monitor/agent.env
 if [[ ! -f "$ENVFILE" ]];then
  if [[ -z "${MONITOR_SERVER:-}" || -z "${MONITOR_AGENT_TOKEN:-}" || -z "${MONITOR_NODE_NAME:-}" ]];then
   echo "Set MONITOR_SERVER=https://your-domain and MONITOR_AGENT_TOKEN=<secret> before first install" >&2
   exit 1
  fi
  case "$MONITOR_SERVER" in https://*|http://127.0.0.1:*) ;; *) echo "MONITOR_SERVER needs HTTPS" >&2;exit 1;; esac
  umask 077
  printf 'MONITOR_SERVER=%s\nMONITOR_AGENT_TOKEN=%s\nMONITOR_NODE_NAME=%s\n' "$MONITOR_SERVER" "$MONITOR_AGENT_TOKEN" "$MONITOR_NODE_NAME" > "$ENVFILE"
 fi
 USERNAME=monitor-agent
 CMD='/usr/local/bin/monitor -mode agent'
fi
chmod 600 "$ENVFILE"
if [[ -e "/etc/systemd/system/monitor-$MODE.service" ]]; then
 echo "Existing systemd unit detected; preserving current port and settings."
 systemctl daemon-reload
 systemctl restart "monitor-$MODE.service"
 systemctl status "monitor-$MODE.service" --no-pager
 exit 0
fi
cat > "/etc/systemd/system/monitor-$MODE.service" <<EOF
[Unit]
Description=Monitor $MODE
Wants=network-online.target
After=network-online.target
[Service]
User=$USERNAME
EnvironmentFile=$ENVFILE
ExecStart=$CMD
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
$( [[ "$MODE" == server ]] && echo StateDirectory=monitor || true )
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable --now "monitor-$MODE.service"
systemctl status "monitor-$MODE.service" --no-pager
