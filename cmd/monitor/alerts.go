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
 }{if _,e:=db.Exec(q);e!=nil{return e}};return nil
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
 rows,e:=db.QueryContext(ctx,"SELECT id,node,kind,start,end FROM alert_events ORDER BY id DESC LIMIT 100");if e!=nil{return out,e};defer rows.Close()
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
func startAlertLoop(s *Store){
 if e:=alertSchema(s.db);e!=nil{return}
 states:=map[string]conditionState{}
 ticker:=time.NewTicker(10*time.Second);defer ticker.Stop()
 for range ticker.C {
  cfg:=getAlertSettings(s.db)
  s.RLock();nodes:=make([]Node,0,len(s.Nodes));for _,n:=range s.Nodes{nodes=append(nodes,n)};s.RUnlock()
  now:=time.Now()
  for _,n:=range nodes{
   if n.LastSeen.IsZero(){continue}
   tests:=[]struct{kind string;bad bool;duration time.Duration}{
    {"offline",now.Sub(n.LastSeen)>=time.Duration(cfg.OfflineSeconds)*time.Second,0},
    {"cpu",n.CPU>=cfg.CPUThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
    {"memory",n.Memory>=cfg.MemoryThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
    {"disk",n.Disk>=cfg.DiskThreshold&&now.Sub(n.LastSeen)<30*time.Second,time.Duration(cfg.DurationSeconds)*time.Second},
   }
   for _,t:=range tests{
    id:=n.Name+"|"+t.kind;st:=states[id]
    if t.bad{
     if st.since.IsZero(){st.since=now}
     if !st.active&&now.Sub(st.since)>=t.duration{
      if _,e:=s.db.Exec("INSERT INTO alert_events(node,kind,start) VALUES(?,?,?)",n.Name,t.kind,now.Unix());e==nil{st.active=true;go deliverWebhook(cfg.Webhook,n.Name,t.kind,"firing")}
     }
    } else {
     st.since=time.Time{}
     if st.active{
      if _,e:=s.db.Exec("UPDATE alert_events SET end=? WHERE node=? AND kind=? AND end=0",now.Unix(),n.Name,t.kind);e==nil{go deliverWebhook(cfg.Webhook,n.Name,t.kind,"resolved");st.active=false}
     }
    }
    states[id]=st
   }
  }
 }
}
