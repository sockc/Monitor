#!/usr/bin/env bash
set -euo pipefail
[[ "${EUID}" -eq 0 ]] || { echo "root required";exit 1; }
DB=/var/lib/monitor/monitor.db
[[ -f "$DB" ]] || { echo "database not found" >&2;exit 1; }
DEST="${1:-/var/backups/monitor}"
mkdir -p "$DEST";chmod 700 "$DEST"
STAMP="$(date -u +%Y%m%d-%H%M%S)"
OUT="$DEST/monitor-$STAMP.db"
# SQLite online backup via Python stdlib sqlite3; no sqlite3 CLI required.
python3 - "$DB" "$OUT" <<'PY'
import sqlite3,sys
src=sqlite3.connect("file:"+sys.argv[1]+"?mode=ro",uri=True)
dst=sqlite3.connect(sys.argv[2])
with dst: src.backup(dst)
dst.close();src.close()
PY
chmod 600 "$OUT"
echo "Backup saved to $OUT"
