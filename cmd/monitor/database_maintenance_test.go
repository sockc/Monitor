package main

import (
 "context"
 "path/filepath"
 "testing"
 "time"
)

func TestWALCheckpointPreservesSamples(t *testing.T){
 path:=filepath.Join(t.TempDir(),"monitor.db")
 db,err:=openDB(path);if err!=nil{t.Fatal(err)};defer db.Close()
 now:=time.Now().UTC()
 for i:=0;i<200;i++{
  s:=Sample{Name:"test-node",Timestamp:now.Add(time.Duration(i)*time.Second),CPU:42}
  if err:=recordSample(db,s);err!=nil{t.Fatal(err)}
 }
 before,err:=walSize(path);if err!=nil{t.Fatal(err)}
 if before==0{t.Fatal("expected WAL after writes")}
 if err:=maybeCheckpointWAL(context.Background(),db,path,1);err!=nil{t.Fatal(err)}
 after,err:=walSize(path);if err!=nil{t.Fatal(err)}
 if after!=0{t.Fatalf("WAL was not truncated: %d bytes",after)}
 var count int
 if err:=db.QueryRow("SELECT COUNT(*) FROM samples").Scan(&count);err!=nil{t.Fatal(err)}
 if count!=200{t.Fatalf("lost samples after checkpoint: %d",count)}
}
func TestRetentionKeepsRecentSamples(t *testing.T){
 db,err:=openDB(filepath.Join(t.TempDir(),"monitor.db"));if err!=nil{t.Fatal(err)};defer db.Close()
 now:=time.Now().UTC()
 if err:=recordSample(db,Sample{Name:"old",Timestamp:now.Add(-31*24*time.Hour)});err!=nil{t.Fatal(err)}
 if err:=recordSample(db,Sample{Name:"new",Timestamp:now.Add(-time.Hour)});err!=nil{t.Fatal(err)}
 removed,err:=pruneOldSamples(context.Background(),db,now);if err!=nil{t.Fatal(err)}
 if removed!=1{t.Fatalf("removed %d samples; want 1",removed)}
 var name string
 if err:=db.QueryRow("SELECT name FROM samples").Scan(&name);err!=nil{t.Fatal(err)}
 if name!="new"{t.Fatalf("retention kept wrong sample: %s",name)}
}
