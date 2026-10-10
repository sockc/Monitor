package main

import (
 "database/sql"
 "path/filepath"
 "testing"
)

func TestNodeSortOrderMigrationPreservesNames(t *testing.T) {
 path := filepath.Join(t.TempDir(), "monitor.db")
 old, err := sql.Open("sqlite", path)
 if err != nil { t.Fatal(err) }
 if _,err=old.Exec("CREATE TABLE node_metadata(name TEXT PRIMARY KEY, display_name TEXT NOT NULL DEFAULT '', group_name TEXT NOT NULL DEFAULT '', notes TEXT NOT NULL DEFAULT '', location TEXT NOT NULL DEFAULT '')");err!=nil{t.Fatal(err)}
 if _,err=old.Exec("INSERT INTO node_metadata(name,display_name,group_name,notes,location) VALUES('n-123','甲骨文ARM','香港','','香港')");err!=nil{t.Fatal(err)}
 if err=old.Close();err!=nil{t.Fatal(err)}
 db,err:=openDB(path);if err!=nil{t.Fatal(err)};defer db.Close()
 var name, group string; var order int
 if err=db.QueryRow("SELECT display_name,group_name,sort_order FROM node_metadata WHERE name='n-123'").Scan(&name,&group,&order);err!=nil{t.Fatal(err)}
 if name!="甲骨文ARM"||group!="香港"||order!=0{t.Fatalf("migration lost node metadata: name=%q group=%q order=%d",name,group,order)}
 if _,err=db.Exec("UPDATE node_metadata SET sort_order=2 WHERE name='n-123'");err!=nil{t.Fatal(err)}
 if err=db.QueryRow("SELECT sort_order FROM node_metadata WHERE name='n-123'").Scan(&order);err!=nil{t.Fatal(err)}
 if order!=2{t.Fatalf("sort order was not saved: %d",order)}
}
