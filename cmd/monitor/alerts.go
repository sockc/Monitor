package main

import (
 "bytes"
 "context"
 "database/sql"
 "encoding/json"
 "errors"
 "net"
 "net/http"
 "net/url"
  "strings"
 "time"
)
type AlertSettings struct {
 OfflineSeconds int `json:"offline_seconds"`
 CPUThreshold float64 `json:"cpu_threshold"`
 MemoryThreshold float64 `json:"memory_threshold"`
 DiskThreshold float64 `json:"disk_threshold"`
 DurationSeconds int `json:"duration_seconds"`
 Webhook string `json:"webhook"`
}
type AlertEvent struct {
 ID int64 `json:"id"`
 Node string `json:"node"`
 Kind string `json:"kind"`
 Start int64 `json:"start"`
 End int64 `json:"end"`
}
type AlertResponse struct{Settings AlertSettings `json:"settings"`; Events []AlertEvent `json:"events"`}
func alertSchema(db *sql.DB)error{
 for _,q:=range []string{
 "CREATE TABLE IF NOT EXISTS alert_settings (id INTEGER PRIMARY KEY CHECK(id=1), config TEXT NOT NULL)",
 "CREATE TABLE IF NOT EXISTS alert_events (id INTEGER PRIMARY KEY AUTOINCREMENT,node TEXT NOT NULL,kind TEXT NOT NULL,start INTEGER NOT NULL,end INTEGER NOT NULL DEFAULT 0)",
 "CREATE INDEX IF NOT EXISTS alert_events_open ON alert_events(node,kind,end)",
 // Older databases may contain duplicate open alarms. Resolve them before
 // adding the partial unique index; retain the earliest active event.
 "UPDATE alert_events SET end=CASE WHEN start>0 THEN start ELSE CAST(strftime('%s','now') AS INTEGER) END WHERE end=0 AND id NOT IN (SELECT MIN(id) FROM alert_events WHERE end=0 GROUP BY node,kind)",
 "CREATE UNIQUE INDEX IF NOT EXISTS alert_one_open ON alert_events(node,kind) WHERE end=0",
 }{if _,e:=db.Exec(q);e!=nil{return e}};return nil
}
// purgeRevokedAlerts is a one-time, idempotent migration for versions which
// removed a node without deleting its alert history. REVOKED tombstones are
// authoritative and ensure no live node's history is touched.
func purgeRevokedAlerts(ctx context.Context,db *sql.DB)(int64,error){
 result,err:=db.ExecContext(ctx,"DELETE FROM alert_events WHERE EXISTS (SELECT 1 FROM agent_tokens WHERE agent_tokens.name=alert_events.node AND agent_tokens.hash='REVOKED')")
 if err!=nil{return 0,err}
 return result.RowsAffected()
}

// insertActiveAlert prevents an in-flight alert-loop snapshot from reopening a
// node's alert after removal. The conditional insert and the deletion
// transaction are serialized by SQLite's write lock.
func insertActiveAlert(db *sql.DB,node,kind string,now int64)(bool,error){
 result,err:=db.Exec("INSERT INTO alert_events(node,kind,start) SELECT ?,?,? WHERE NOT EXISTS (SELECT 1 FROM agent_tokens WHERE name=? AND hash='REVOKED') ON CONFLICT DO NOTHING",node,kind,now,node)
 if err!=nil{return false,err}
 n,err:=result.RowsAffected()
 return n==1,err
}

func defaultAlertSettings()AlertSettings{return AlertSettings{OfflineSeconds:60,CPUThreshold:90,MemoryThreshold:90,DiskThreshold:90,DurationSeconds:300}}
func getAlertSettings(db *sql.DB)AlertSettings{
 cfg:=defaultAlertSettings();var raw string
 if db.QueryRow("SELECT config FROM alert_settings WHERE id=1").Scan(&raw)==nil{json.Unmarshal([]byte(raw),&cfg)}
 return cfg
}
func writeAlertSettings(db *sql.DB,c AlertSettings)error{raw,e:=json.Marshal(c);if e!=nil{return e};_,e=db.Exec("INSERT INTO alert_settings(id,config) VALUES(1,?) ON CONFLICT(id) DO UPDATE SET config=excluded.config",string(raw));return e}
func readAlerts(ctx context.Context,db *sql.DB)(AlertResponse,error){
 out:=AlertResponse{Settings:getAlertSettings(db),Events:[]AlertEvent{}}
 rows,e:=db.QueryContext(ctx,"SELECT a.id,a.node,a.kind,a.start,a.end FROM alert_events a WHERE NOT EXISTS (SELECT 1 FROM agent_tokens t WHERE t.name=a.node AND t.hash='REVOKED') ORDER BY a.id DESC LIMIT 100");if e!=nil{return out,e};defer rows.Close()
 for rows.Next(){var x AlertEvent;if e=rows.Scan(&x.ID,&x.Node,&x.Kind,&x.Start,&x.End);e!=nil{return out,e};out.Events=append(out.Events,x)}
 return out,rows.Err()
}
func validWebhook(raw string)bool{
 u,e:=url.Parse(raw);if e!=nil||u.Scheme!="https"||u.User!=nil||u.Hostname()==""||u.Fragment!=""{return false}
 // Block literal loopback/private destinations. DNS is re-checked at delivery time.
 if ip:=net.ParseIP(u.Hostname());ip!=nil{return ip.IsGlobalUnicast()&&!ip.IsPrivate()}
 return !strings.EqualFold(u.Hostname(),"localhost")&&!strings.HasSuffix(strings.ToLower(u.Hostname()),".local")
}
func deliverWebhook(target,node,kind,status string){
 if target==""{return}
 ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel()
 b,_:=json.Marshal(map[string]string{"node":node,"event":kind,"status":status})
 req,e:=http.NewRequestWithContext(ctx,"POST",target,bytes.NewReader(b));if e!=nil{return}
 req.Header.Set("Content-Type","application/json")
 // Outbound webhook is opt-in; redirect disabled.
 client:=&http.Client{Timeout:5*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return errors.New("redirect forbidden")}}
 res,e:=client.Do(req);if e==nil{res.Body.Close()}
}
type conditionState struct{since time.Time;active bool}
func expirationKind(expires,zone string,now time.Time)string{
 if expires==""{return ""};loc,e:=time.LoadLocation(zone);if e!=nil{loc=time.UTC}
 end,e:=time.ParseInLocation("2006-01-02",expires,loc);if e!=nil{return ""}
 local:=now.In(loc);today:=time.Date(local.Year(),local.Month(),local.Day(),0,0,0,0,loc)
 if today.After(end){return "expiry_overdue"}
 if !today.Before(end.AddDate(0,0,-7)){return "expiry_7"}
 if !today.Before(end.AddDate(0,0,-15)){return "expiry_15"}
 if !today.Before(end.AddDate(0,0,-30)){return "expiry_30"}
 return ""
}
func quotaKind(p PeriodTraffic)string{
 if p.QuotaGB<=0||!p.HasSamples{return ""}
 quota:=p.QuotaGB*1000000000;used:=float64(p.MonthRX)+float64(p.MonthTX)
 switch{case used>=quota:return "quota_100";case used>=quota*0.9:return "quota_90";case used>=quota*0.8:return "quota_80";default:return ""}
}

func startAlertLoop(s *Store){
 if e:=alertSchema(s.db);e!=nil{return}
 states:=map[string]conditionState{}
 ticker:=time.NewTicker(10*time.Second);defer ticker.Stop()
 for range ticker.C {
  cfg:=getAlertSettings(s.db)
  s.RLock();nodes:=make([]Node,0,len(s.Nodes));for _,n:=range s.Nodes{nodes=append(nodes,n)};s.RUnlock()
  now:=time.Now()
  summaries,e:=periodTraffic(context.Background(),s.db);if e!=nil{continue}
  for _,n:=range nodes{
   if n.LastSeen.IsZero(){continue}
   tests:=[]struct{kind string;bad bool;duration time.Duration}{
    {"offline",now.Sub(n.LastSeen)>=time.Duration(cfg.OfflineSeconds)*time.Second,0},
    {"cpu",n.CPU>=cfg.CPUThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
    {"memory",n.Memory>=cfg.MemoryThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
    {"disk",n.Disk>=cfg.DiskThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
    // Three TCP connect attempts to the Monitor HTTPS endpoint. This is a
    // connectivity failure ratio, not ICMP loss; only fresh probes qualify.
    {"network",n.NetworkFailurePercent!=nil&&n.NetworkCheckedAt>0&&now.Unix()-n.NetworkCheckedAt<150&&now.Sub(n.LastSeen)<30*time.Second&&*n.NetworkFailurePercent>=66,120*time.Second},
   }
   if p,ok:=summaries[n.Name];ok{
    stage:=quotaKind(p);for _,kind:=range []string{"quota_80","quota_90","quota_100"}{tests=append(tests,struct{kind string;bad bool;duration time.Duration}{kind,kind==stage,0})}
    expiry:=expirationKind(p.ExpiresOn,p.Timezone,now);for _,kind:=range []string{"expiry_30","expiry_15","expiry_7","expiry_overdue"}{tests=append(tests,struct{kind string;bad bool;duration time.Duration}{kind,kind==expiry,0})}
   }
   for _,t:=range tests{
    id:=n.Name+"|"+t.kind;st,seen:=states[id]
    if !seen{var count int;_ = s.db.QueryRow("SELECT COUNT(*) FROM alert_events WHERE node=? AND kind=? AND end=0",n.Name,t.kind).Scan(&count);st.active=count>0}
    if t.bad{
     if st.since.IsZero(){st.since=now}
     if !st.active&&now.Sub(st.since)>=t.duration{
      if inserted,e:=insertActiveAlert(s.db,n.Name,t.kind,now.Unix());e==nil&&inserted{st.active=true;go deliverWebhook(cfg.Webhook,n.Name,t.kind,"firing")}
     }
    } else {
     st.since=time.Time{}
     if st.active{
      if result,e:=s.db.Exec("UPDATE alert_events SET end=? WHERE node=? AND kind=? AND end=0",now.Unix(),n.Name,t.kind);e==nil{if changed,_:=result.RowsAffected();changed>0{go deliverWebhook(cfg.Webhook,n.Name,t.kind,"resolved")};st.active=false}
     }
    }
    states[id]=st
   }
  }
 }
}
