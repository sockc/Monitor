#!/usr/bin/env bash
set -euo pipefail
[[ "${EUID}" -eq 0 ]] || { echo "root required";exit 1; }
MODE="${1:-}";[[ "$MODE" == server || "$MODE" == agent ]] || { echo "Usage: sudo bash upgrade.sh server|agent";exit 1; }
VERSION="${MONITOR_VERSION:-v0.9.15}"
ARCH="$(uname -m)";case "$ARCH" in x86_64)ARCH=amd64;;aarch64|arm64)ARCH=arm64;;*)exit 1;;esac
ROOT="https://github.com/sockc/Monitor/releases/download/$VERSION"
T="$(mktemp -d)";trap 'rm -rf "$T"' EXIT
curl -fLsS --retry 3 "$ROOT/monitor-linux-$ARCH" -o "$T/monitor"
curl -fLsS --retry 3 "$ROOT/SHA256SUMS" -o "$T/SHA256SUMS"
(cd "$T";grep "monitor-linux-$ARCH$" SHA256SUMS | sed -E "s@[^[:space:]]*monitor-linux-$ARCH@monitor@" | sha256sum -c -)
systemctl is-active --quiet "monitor-$MODE.service" || { echo "service not running" >&2;exit 1; }
if [[ "$MODE" == server ]];then
 SCRIPT="$(dirname "$0")/backup.sh"
 if [[ -f "$SCRIPT" ]];then bash "$SCRIPT";else echo "Run a database backup before upgrading" >&2;exit 1;fi
fi
cp -p /usr/local/bin/monitor "$T/previous"
install -m 0755 "$T/monitor" /usr/local/bin/monitor.new
mv -f /usr/local/bin/monitor.new /usr/local/bin/monitor
systemctl restart "monitor-$MODE.service"
sleep 3
HEALTH_OK=true
if [[ "$MODE" == server ]]; then
 PORT="${MONITOR_HEALTH_PORT:-}"
 if [[ -z "$PORT" ]]; then
   PORT="$(systemctl cat monitor-server.service | sed -nE 's/.*-listen 127[.]0[.]0[.]1:([0-9]+).*/\1/p' | tail -1)"
 fi
 if [[ -z "$PORT" && -f /etc/monitor/server.env ]]; then
   PORT="$(sed -nE 's/^MONITOR_LISTEN=127[.]0[.]0[.]1:([0-9]+)$/\1/p' /etc/monitor/server.env | tail -1)"
 fi
 PORT="${PORT:-8090}"
 curl -fsS --max-time 3 "http://127.0.0.1:$PORT/healthz" >/dev/null || HEALTH_OK=false
fi
if ! systemctl is-active --quiet "monitor-$MODE.service" || [[ "$HEALTH_OK" != true ]]; then
 install -m 0755 "$T/previous" /usr/local/bin/monitor.rollback
 mv -f /usr/local/bin/monitor.rollback /usr/local/bin/monitor
 systemctl restart "monitor-$MODE.service"
 echo "Upgrade failed; binary rolled back" >&2
 exit 1
fi
echo "Updated to $VERSION"
