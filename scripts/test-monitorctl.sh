#!/usr/bin/env bash
# Non-destructive smoke checks for the lifecycle wrapper.
set -euo pipefail
cd "$(dirname "$0")/.."

bash -n scripts/monitorctl scripts/install-release.sh scripts/upgrade.sh scripts/backup.sh
help_output="$(bash scripts/monitorctl --help)"
[[ "$help_output" == *"install server"* ]]
[[ "$help_output" == *"install agent"* ]]
[[ "$help_output" == *"uninstall server|agent"* ]]
[[ "$help_output" == *"self-update"* ]]

version_output="$(MONITOR_VERSION=v9-test bash scripts/monitorctl version)"
[[ "$version_output" == *"v9-test"* ]]
if bash scripts/monitorctl install invalid >/dev/null 2>&1;then
 echo "invalid role was accepted" >&2
 exit 1
fi
if bash scripts/monitorctl unsupported >/dev/null 2>&1;then
 echo "unsupported command was accepted" >&2
 exit 1
fi

# This regression previously created server.env with MONITOR_LISTEN alone,
# preventing generation of the required secrets on a clean installation.
grep -q 'MONITOR_AGENT_TOKEN=%s' scripts/install-release.sh
grep -q 'MONITOR_ADMIN_TOKEN=%s' scripts/install-release.sh
grep -q 'MONITOR_LISTEN=127.0.0.1:%s' scripts/install-release.sh
echo "monitorctl command and installation smoke checks passed"
