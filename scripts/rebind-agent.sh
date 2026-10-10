#!/usr/bin/env bash
# Reconnect an EXISTING agent to an immutable Monitor node ID. Does not
# reinstall the Agent binary, delete history or replace unrelated settings.
set -euo pipefail
[[ ${EUID} -eq 0 ]] || { echo "请使用 sudo 运行" >&2;exit 1; }
: "${MONITOR_SERVER:?Missing Server URL}"
: "${MONITOR_NODE_NAME:?Missing generated node ID}"
: "${MONITOR_AGENT_TOKEN:?Missing node token}"
[[ "$MONITOR_NODE_NAME" =~ ^[a-zA-Z0-9_-]{1,64}$ ]] || { echo "Invalid node ID" >&2;exit 1; }
[[ "$MONITOR_AGENT_TOKEN" =~ ^[a-fA-F0-9]{64}$ ]] || { echo "Invalid node token" >&2;exit 1; }
[[ "$MONITOR_SERVER" =~ ^https://[a-zA-Z0-9._-]+(:[0-9]{1,5})?(/)?$ ]] || { echo "Server HTTPS URL is invalid" >&2;exit 1; }
systemctl cat monitor-agent.service >/dev/null 2>&1 || { echo "Agent is not installed" >&2;exit 1; }
FILE="/etc/monitor/agent.env"
[[ -f "$FILE" ]] || { echo "Agent configuration not found" >&2;exit 1; }
umask 077
tmp="$(mktemp /etc/monitor/.agent.env.XXXXXX)"
trap 'rm -f "$tmp"' EXIT
seen_server=0;seen_name=0;seen_token=0
while IFS= read -r line || [[ -n "$line" ]];do
 case "$line" in
 MONITOR_SERVER=*)printf 'MONITOR_SERVER=%s\n' "$MONITOR_SERVER" >> "$tmp";seen_server=1;;
 MONITOR_NODE_NAME=*)printf 'MONITOR_NODE_NAME=%s\n' "$MONITOR_NODE_NAME" >> "$tmp";seen_name=1;;
 MONITOR_AGENT_TOKEN=*)printf 'MONITOR_AGENT_TOKEN=%s\n' "$MONITOR_AGENT_TOKEN" >> "$tmp";seen_token=1;;
 *)printf '%s\n' "$line" >> "$tmp";;
 esac
done < "$FILE"
if (( !seen_server ));then printf 'MONITOR_SERVER=%s\n' "$MONITOR_SERVER" >> "$tmp";fi
if (( !seen_name ));then printf 'MONITOR_NODE_NAME=%s\n' "$MONITOR_NODE_NAME" >> "$tmp";fi
if (( !seen_token ));then printf 'MONITOR_AGENT_TOKEN=%s\n' "$MONITOR_AGENT_TOKEN" >> "$tmp";fi
chmod 600 "$tmp"
chown --reference="$FILE" "$tmp"
mv -f "$tmp" "$FILE"
systemctl restart monitor-agent.service
sleep 2
systemctl is-active --quiet monitor-agent.service || { echo "Agent restart failed. Check journalctl -u monitor-agent" >&2;exit 1; }
echo "Agent 已重新绑定节点 $MONITOR_NODE_NAME；请等待上报刷新。"
