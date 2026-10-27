#!/usr/bin/env bash
set -euo pipefail
[[ "${EUID}" -eq 0 ]] || { echo "root required" >&2; exit 1; }
DB=/var/lib/monitor/monitor.db
[[ -f "$DB" ]] || { echo "database not found" >&2; exit 1; }
DEST="${1:-/var/backups/monitor}"
mkdir -p "$DEST"; chmod 700 "$DEST"
STAMP="$(date -u +%Y%m%d-%H%M%S)"
OUT="$DEST/monitor-$STAMP.db"
# SQLite online backup includes committed WAL transactions. Verify the
# copy before publishing, and remove a partial file on any error.
python3 - "$DB" "$OUT" <<'PY'
import os,sqlite3,sys
database,output=sys.argv[1:]
partial=output+".partial-"+str(os.getpid())
src=None;dst=None
try:
    src=sqlite3.connect("file:"+database+"?mode=ro",uri=True,timeout=60)
    dst=sqlite3.connect(partial,timeout=60)
    with dst:
        src.backup(dst,pages=256,sleep=0.1)
    if dst.execute("PRAGMA quick_check").fetchone()[0]!="ok":
        raise RuntimeError("backup quick_check failed")
    dst.close();dst=None
    src.close();src=None
    os.chmod(partial,0o600)
    os.replace(partial,output)
    print("Backup saved to",output)
except Exception:
    if dst: dst.close()
    if src: src.close()
    if os.path.exists(partial): os.unlink(partial)
    raise
PY
