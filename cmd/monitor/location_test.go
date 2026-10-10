package main

import (
 "database/sql"
 "path/filepath"
 "testing"
 _ "modernc.org/sqlite"
)

func TestLocationMigrationPreservesExistingMetadata(t *testing.T){
 path:=filepath.Join(t.TempDir(),"existing.db")
 old,e:=sql.Open("sqlite",path);if e!=nil{t.Fatal(e)}
 if _,e=old.Exec("CREATE TABLE node_metadata(name TEXT PRIMARY KEY,display_name TEXT NOT NULL DEFAULT '',group_name TEXT NOT NULL DEFAULT '',notes TEXT NOT NULL DEFAULT '')");e!=nil{t.Fatal(e)}
 if _,e=old.Exec("INSERT INTO node_metadata(name,display_name,group_name,notes) VALUES('hk-01','香港机','生产','保留备注')");e!=nil{t.Fatal(e)}
 if e=old.Close();e!=nil{t.Fatal(e)}
 db,e:=openDB(path);if e!=nil{t.Fatal(e)}
 var name,group,location,notes string
 if e=db.QueryRow("SELECT display_name,group_name,location,notes FROM node_metadata WHERE name='hk-01'").Scan(&name,&group,&location,&notes);e!=nil{t.Fatal(e)}
 if name!="香港机"||group!="生产"||location!=""||notes!="保留备注"{t.Fatalf("migration changed metadata: %q, %q, %q, %q",name,group,location,notes)}
 if _,e=db.Exec("UPDATE node_metadata SET location='中国香港 · 沙田' WHERE name='hk-01'");e!=nil{t.Fatal(e)}
 if e=db.Close();e!=nil{t.Fatal(e)}
 // Re-running migrations on an already upgraded database must succeed.
 db,e=openDB(path);if e!=nil{t.Fatal(e)};defer db.Close()
 if e=db.QueryRow("SELECT location FROM node_metadata WHERE name='hk-01'").Scan(&location);e!=nil{t.Fatal(e)}
 if location!="中国香港 · 沙田"{t.Fatalf("persisted location lost: %q",location)}
}
