#!/usr/bin/env bash
set -euo pipefail
[[ "${EUID}" -eq 0 ]] || { echo "root required" >&2;exit 1; }
[[ "${MONITOR_BACKUP_CONFIRMED:-}" == 1 ]] || { echo "Run via sudo monitorctl db-maintain: verified backup required" >&2; exit 1; }
command -v lsof >/dev/null || { echo "lsof required to verify database is offline" >&2;exit 1; }
systemctl is-active --quiet monitor-server.service || { echo "Monitor Server must be running" >&2;exit 1; }
DB=/var/lib/monitor/monitor.db
systemctl stop monitor-server.service
trap 'systemctl start monitor-server.service' EXIT
if lsof "$DB" "$DB-wal" "$DB-shm" 2>/dev/null;then
 echo "Database still open by another process; checkpoint refused" >&2
 exit 1
fi
runuser -u monitor -- python3 - "$DB" <<'PY'
import sqlite3,sys
db=sqlite3.connect(sys.argv[1],timeout=120)
try:
    print("SQLite mode:",db.execute("PRAGMA journal_mode").fetchone()[0])
    result=db.execute("PRAGMA wal_checkpoint(TRUNCATE)").fetchone()
    print("WAL checkpoint (busy,frames,checkpointed):",result)
    if result[0]!=0:
        raise RuntimeError("WAL checkpoint incomplete; database files preserved")
    check=db.execute("PRAGMA quick_check").fetchone()[0]
    print("Database quick_check:",check)
    if check!="ok":
        raise RuntimeError("Database quick_check failed")
finally:
    db.close()
PY
ls -lh /var/lib/monitor/monitor.db*
df -h /
