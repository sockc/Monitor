package main

import (
 "context"
 "testing"
)

func TestRemoveNodeDeletesItsAlertsOnly(t *testing.T) {
 db:=testUptimeDB(t)
 if err:=alertSchema(db);err!=nil{t.Fatal(err)}
 for _,node:=range []string{"deleted","active"}{
  if _,err:=db.Exec("INSERT INTO samples(name,ts,cpu,memory,disk,payload) VALUES(?,?,?,?,?,?)",node,100,0,0,0,"{}");err!=nil{t.Fatal(err)}
  if _,err:=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES(?,?,?,?)",node,"offline",100,0);err!=nil{t.Fatal(err)}
  if _,err:=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES(?,?,?,?)",node,"cpu",80,90);err!=nil{t.Fatal(err)}
 }
 if err:=removeNode(context.Background(),db,"deleted");err!=nil{t.Fatal(err)}
 var count int
 if err:=db.QueryRow("SELECT COUNT(*) FROM alert_events WHERE node='deleted'").Scan(&count);err!=nil||count!=0{t.Fatalf("deleted alerts remain: %d, %v",count,err)}
 if err:=db.QueryRow("SELECT COUNT(*) FROM alert_events WHERE node='active'").Scan(&count);err!=nil||count!=2{t.Fatalf("active alerts changed: %d, %v",count,err)}
 var tombstone string
 if err:=db.QueryRow("SELECT hash FROM agent_tokens WHERE name='deleted'").Scan(&tombstone);err!=nil||tombstone!="REVOKED"{t.Fatalf("missing tombstone: %s, %v",tombstone,err)}
 result,err:=readAlerts(context.Background(),db)
 if err!=nil||len(result.Events)!=2{t.Fatalf("unexpected alert list: %+v, %v",result.Events,err)}
 for _,event:=range result.Events{if event.Node!="active"{t.Fatalf("deleted node shown: %+v",event)}}
}

func TestRevokeNodeRemovesAlertButKeepsSamples(t *testing.T){
 db:=testUptimeDB(t)
 if err:=alertSchema(db);err!=nil{t.Fatal(err)}
 _,err:=db.Exec("INSERT INTO samples(name,ts,cpu,memory,disk,payload) VALUES('revoked',100,0,0,0,'{}')")
 if err!=nil{t.Fatal(err)}
 _,err=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES('revoked','offline',100,0)")
 if err!=nil{t.Fatal(err)}
 if err:=revokeNode(context.Background(),db,"revoked");err!=nil{t.Fatal(err)}
 var alerts,samples int
 db.QueryRow("SELECT COUNT(*) FROM alert_events WHERE node='revoked'").Scan(&alerts)
 db.QueryRow("SELECT COUNT(*) FROM samples WHERE name='revoked'").Scan(&samples)
 if alerts!=0||samples!=1{t.Fatalf("revoke alerts=%d samples=%d",alerts,samples)}
}

func TestLegacyDeletedAlertsCleanupAndLateInsert(t *testing.T){
 db:=testUptimeDB(t)
 if err:=alertSchema(db);err!=nil{t.Fatal(err)}
 _,err:=db.Exec("INSERT INTO agent_tokens(name,hash) VALUES('old-node','REVOKED')")
 if err!=nil{t.Fatal(err)}
 for _,node:=range []string{"old-node","live-node"}{
  for _,event:=range []struct{kind string;end int}{{"offline",0},{"memory",300}}{
   if _,err=db.Exec("INSERT INTO alert_events(node,kind,start,end) VALUES(?,?,?,?)",node,event.kind,100,event.end);err!=nil{t.Fatal(err)}
  }
 }
 // Even before the migration, the API must not leak revoked events.
 before,err:=readAlerts(context.Background(),db)
 if err!=nil||len(before.Events)!=2{t.Fatalf("old alerts leaked before cleanup: %+v %v",before.Events,err)}
 removed,err:=purgeRevokedAlerts(context.Background(),db)
 if err!=nil||removed!=2{t.Fatalf("cleanup removed=%d, err=%v",removed,err)}
 repeat,err:=purgeRevokedAlerts(context.Background(),db)
 if err!=nil||repeat!=0{t.Fatalf("cleanup is not idempotent: %d, %v",repeat,err)}
 inserted,err:=insertActiveAlert(db,"old-node","offline",400)
 if err!=nil||inserted{t.Fatalf("reopened deleted node: inserted=%v err=%v",inserted,err)}
 inserted,err=insertActiveAlert(db,"live-node","network",400)
 if err!=nil||!inserted{t.Fatalf("live node alert blocked: %v %v",inserted,err)}
 after,err:=readAlerts(context.Background(),db)
 if err!=nil||len(after.Events)!=3{t.Fatalf("live events altered: %+v %v",after.Events,err)}
}
