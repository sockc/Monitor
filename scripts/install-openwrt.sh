#!/bin/sh
set -eu
# Installation: sh install-openwrt.sh https://monitor.example.com NODE_NAME
SERVER="${1:-}"; NAME="${2:-}"
[ -n "$SERVER" ] && [ -n "$NAME" ] || { echo 'Usage: sh install-openwrt.sh https://server NODE_NAME'; exit 1; }
if [ ! -e /etc/monitor-agent.json ]; then
 printf 'Paste Agent token: ' >&2
 stty -echo 2>/dev/null || true
 IFS= read -r TOKEN
 stty echo 2>/dev/null || true
 printf '\n' >&2
 [ -n "$TOKEN" ] || { echo 'Token required'; exit 1; }
else TOKEN=''; fi
[ "$(uname -m)" = "aarch64" ] || { echo "Only aarch64 is supported"; exit 1; }
[ -f /lib/functions/procd.sh ] || { echo "procd is required"; exit 1; }
case "$SERVER" in https://*) ;; *) echo "HTTPS server URL required"; exit 1;; esac
[ -n "$NAME" ] && [ -n "$TOKEN" ] || exit 1
RELEASE="v0.9.16.1"
BASE="https://github.com/sockc/Monitor/releases/download/$RELEASE"
BINARY="monitor-openwrt-linux-arm64"
TMP="$(mktemp -d /tmp/monitor-install.XXXXXX)"
trap 'rm -rf "$TMP"' EXIT HUP INT TERM
if command -v curl >/dev/null 2>&1; then
 curl -fL "$BASE/$BINARY" -o "$TMP/$BINARY"
 curl -fL "$BASE/SHA256SUMS" -o "$TMP/SHA256SUMS"
else
 wget -O "$TMP/$BINARY" "$BASE/$BINARY"
 wget -O "$TMP/SHA256SUMS" "$BASE/SHA256SUMS"
fi
(cd "$TMP" && grep "  $BINARY$" SHA256SUMS | sha256sum -c -)
BYTES="$(wc -c < "$TMP/$BINARY")"
FREE="$(df -k /overlay | awk 'NR==2 {print $4}')"
NEEDED="$(( (BYTES+1023)/1024 + 10240 ))"
[ -n "$FREE" ] && [ "$FREE" -gt "$NEEDED" ] || { echo "Need at least 10 MiB flash headroom"; exit 1; }
umask 077
if [ ! -e /etc/monitor-agent.json ]; then
 escape() { printf '%s' "$1" | sed 's/\\/\\\\/g;s/"/\\"/g'; }
 printf '{"server":"%s","name":"%s","token":"%s","interval_seconds":10}\n' \
   "$(escape "$SERVER")" "$(escape "$NAME")" "$(escape "$TOKEN")" > "$TMP/config"
 cp "$TMP/config" /etc/monitor-agent.json
 chmod 600 /etc/monitor-agent.json
fi
cat > "$TMP/init" <<'EOF'
#!/bin/sh /etc/rc.common
USE_PROCD=1
START=95
STOP=15
start_service() {
 procd_open_instance
 procd_set_param command /usr/bin/monitor-openwrt-agent -config /etc/monitor-agent.json
 procd_set_param respawn 3600 5 5
 procd_set_param stderr 1
 procd_close_instance
}
EOF
[ ! -x /etc/init.d/monitor-agent ] || /etc/init.d/monitor-agent stop || true
cp "$TMP/$BINARY" /usr/bin/monitor-openwrt-agent
chmod 755 /usr/bin/monitor-openwrt-agent
cp "$TMP/init" /etc/init.d/monitor-agent
chmod 755 /etc/init.d/monitor-agent
/etc/init.d/monitor-agent enable
/etc/init.d/monitor-agent start
echo "Monitor OpenWrt Agent started."
