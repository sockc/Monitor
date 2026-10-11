package main

import (
 "context"
 "database/sql"
 "path/filepath"
 "testing"
 "time"
)
func testUptimeDB(t *testing.T)*sql.DB{
 t.Helper()
 db,e:=openDB(filepath.Join(t.TempDir(),"test.db"));if e!=nil{t.Fatal(e)}
 t.Cleanup(func(){db.Close()})
 return db
}
func TestHeartbeatCoverage(t *testing.T){
 db:=testUptimeDB(t)
 now:=time.Unix(1_800_000_000,0).UTC()
 first:=now.Unix()/60-4
 // 5 elapsed minutes; 3 minutes have at least one heartbeat.
 for _,minute:=range []int64{first,first,first+1,first+3}{
  _,err:=db.Exec("INSERT OR REPLACE INTO samples(name,ts,cpu,memory,disk,payload) VALUES(?,?,?,?,?,?)","node",minute*60+10,10,20,30,"{}")
  if err!=nil{t.Fatal(err)}
 }
 windows,err:=queryUptime(context.Background(),db,"node",now)
 if err!=nil{t.Fatal(err)}
 if len(windows)!=3{t.Fatalf("want three windows: %+v",windows)}
 if windows[0].ExpectedMinutes!=5||windows[0].ObservedMinutes!=3||windows[0].Percent==nil||*windows[0].Percent!=60{
  t.Fatalf("unexpected coverage: %+v",windows[0])
 }
 for _,w:=range windows[1:]{if w.Percent==nil||*w.Percent!=60{t.Fatalf("new node should use first heartbeat, got %+v",w)}}
}
func TestCoverageNoSamples(t *testing.T){
 db:=testUptimeDB(t)
 got,err:=queryUptime(context.Background(),db,"not-created",time.Now())
 if err!=nil{t.Fatal(err)}
 for _,w:=range got {if w.Percent!=nil||w.ExpectedMinutes!=0{t.Fatalf("expected unknown coverage: %+v",w)}}
}
func TestAlertSchemaRepairsDuplicateOpen(t *testing.T){
 db:=testUptimeDB(t)
 // Older databases may have two open alerts for the same condition.
 _,err:=db.Exec("CREATE TABLE IF NOT EXISTS alert_events (id INTEGER PRIMARY KEY AUTOINCREMENT,node TEXT NOT NULL,kind TEXT NOT NULL,start INTEGER NOT NULL,end INTEGER NOT NULL DEFAULT 0)")
 if err!=nil{t.Fatal(err)}
 for i:=0;i<2;i++{if _,e:=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES('node','cpu',100,0)");e!=nil{t.Fatal(e)}}
 if err:=alertSchema(db);err!=nil{t.Fatal(err)}
 var active int
 if err:=db.QueryRow("SELECT COUNT(*) FROM alert_events WHERE node='node' AND kind='cpu' AND end=0").Scan(&active);err!=nil{t.Fatal(err)}
 if active!=1{t.Fatalf("got %d open events",active)}
 if _,err=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES('node','cpu',101,0)");err==nil{t.Fatal("duplicate active alert unexpectedly allowed")}
}
