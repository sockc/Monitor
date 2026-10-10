#!/usr/bin/env bash
set -euo pipefail
[[ "${EUID}" -eq 0 ]] || { echo "Please run as root" >&2;exit 1; }
MODE="${1:-}"
[[ "$MODE" == server || "$MODE" == agent ]] || { echo "Usage: sudo bash install-release.sh server|agent" >&2;exit 1; }
command -v systemctl >/dev/null || { echo "systemd is required" >&2;exit 1; }
command -v curl >/dev/null || { echo "curl is required" >&2;exit 1; }
command -v sha256sum >/dev/null || { echo "sha256sum is required" >&2;exit 1; }

# Existing services are managed via upgrade.sh; never silently overwrite a
# shared Server/Agent binary or change credentials during re-install.
if systemctl cat "monitor-$MODE.service" >/dev/null 2>&1;then
 echo "monitor-$MODE.service is already installed. Use monitorctl upgrade $MODE." >&2
 exit 2
fi

ARCH="$(uname -m)"
case "$ARCH" in x86_64) ARCH=amd64;; aarch64|arm64) ARCH=arm64;; *) echo "Unsupported CPU $ARCH" >&2;exit 1;;esac
VERSION="${MONITOR_VERSION:-v0.9.8}"
BASE="https://github.com/sockc/Monitor/releases/download/$VERSION"

# Validate prerequisites and agent credentials BEFORE installing any binary.
if [[ "$MODE" == agent ]];then
 : "${MONITOR_SERVER:?Set MONITOR_SERVER to your HTTPS dashboard URL}"
 : "${MONITOR_AGENT_TOKEN:?Set MONITOR_AGENT_TOKEN to a node-scoped token}"
 : "${MONITOR_NODE_NAME:?Set MONITOR_NODE_NAME to an ASCII node ID}"
 [[ "$MONITOR_NODE_NAME" =~ ^[a-zA-Z0-9_-]{1,64}$ ]] || { echo "Invalid node ID" >&2;exit 1; }
 case "$MONITOR_SERVER" in https://*|http://127.0.0.1:*) ;; *) echo "MONITOR_SERVER requires HTTPS" >&2;exit 1;;esac
 [[ "$MONITOR_SERVER" != *[[:space:]]* && "$MONITOR_AGENT_TOKEN" != *[[:space:]]* ]] || { echo "Whitespace in server URL/token" >&2;exit 1; }
else
 PORT="${MONITOR_PORT:-8090}"
 [[ "$PORT" =~ ^[0-9]{1,5}$ ]] && ((10#$PORT>=1 && 10#$PORT<=65535)) || { echo "Invalid MONITOR_PORT" >&2;exit 1; }
 command -v openssl >/dev/null || { echo "openssl is required" >&2;exit 1; }
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
curl --proto '=https' --tlsv1.2 -fLSs --retry 3 "$BASE/monitor-linux-$ARCH" -o "$TMP/monitor"
curl --proto '=https' --tlsv1.2 -fLSs --retry 3 "$BASE/SHA256SUMS" -o "$TMP/SHA256SUMS"
(cd "$TMP";grep "monitor-linux-$ARCH$" SHA256SUMS | sed -E "s@[^[:space:]]*monitor-linux-$ARCH@monitor@" | sha256sum -c -)

install -d -m 0700 /etc/monitor
umask 077
if [[ "$MODE" == server ]];then
 id monitor >/dev/null 2>&1 || useradd --system --home-dir /var/lib/monitor --shell /usr/sbin/nologin monitor
 install -d -m 0700 -o monitor -g monitor /var/lib/monitor
 ENVFILE=/etc/monitor/server.env
 if [[ ! -e "$ENVFILE" ]];then
  # One atomic initialization of listen address AND both independent secrets.
  printf 'MONITOR_LISTEN=127.0.0.1:%s\nMONITOR_AGENT_TOKEN=%s\nMONITOR_ADMIN_TOKEN=%s\n' \
   "$PORT" "$(openssl rand -hex 32)" "$(openssl rand -hex 32)" > "$TMP/server.env"
  install -m 0600 "$TMP/server.env" "$ENVFILE"
 else
  # Preserve all existing values; only fill missing keys from an incomplete
  # older installation. Do not rotate previously issued credentials.
  if grep -Eq '^MONITOR_LISTEN=.*\\n' "$ENVFILE";then
   echo "Malformed existing server.env: contains literal backslash-n. Repair manually before installing." >&2
   exit 1
  fi
  if ! grep -q '^MONITOR_LISTEN=' "$ENVFILE";then
   printf 'MONITOR_LISTEN=127.0.0.1:%s\n' "$PORT" >> "$ENVFILE"
  fi
  if ! grep -q '^MONITOR_AGENT_TOKEN=' "$ENVFILE";then
   printf 'MONITOR_AGENT_TOKEN=%s\n' "$(openssl rand -hex 32)" >> "$ENVFILE"
  fi
  if ! grep -q '^MONITOR_ADMIN_TOKEN=' "$ENVFILE";then
   printf 'MONITOR_ADMIN_TOKEN=%s\n' "$(openssl rand -hex 32)" >> "$ENVFILE"
  fi
 fi
 USERNAME=monitor
 CMD='/usr/local/bin/monitor -mode server'
else
 id monitor-agent >/dev/null 2>&1 || useradd --system --home-dir /nonexistent --shell /usr/sbin/nologin monitor-agent
 ENVFILE=/etc/monitor/agent.env
 if [[ ! -e "$ENVFILE" ]];then
  printf 'MONITOR_SERVER=%s\nMONITOR_AGENT_TOKEN=%s\nMONITOR_NODE_NAME=%s\n' \
   "$MONITOR_SERVER" "$MONITOR_AGENT_TOKEN" "$MONITOR_NODE_NAME" > "$TMP/agent.env"
  install -m 0600 "$TMP/agent.env" "$ENVFILE"
 fi
 USERNAME=monitor-agent
 CMD='/usr/local/bin/monitor -mode agent'
fi
chmod 600 "$ENVFILE"

install -m 0755 "$TMP/monitor" /usr/local/bin/monitor
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
chmod 644 "/etc/systemd/system/monitor-$MODE.service"
systemctl daemon-reload
systemctl enable --now "monitor-$MODE.service"
systemctl --no-pager --full status "monitor-$MODE.service"
echo "Monitor $MODE installed ($VERSION). Configuration preserved at $ENVFILE."
