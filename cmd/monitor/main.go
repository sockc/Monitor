package main

import (
 "crypto/subtle"
 
 "encoding/json"
 "database/sql"
 "crypto/sha256"
 "encoding/hex"
 "crypto/rand"
 _ "modernc.org/sqlite"
 "flag"
 "fmt"
 "html/template"
 "log"
 "net/http"
 "os"
 "path/filepath"
 "runtime"
 "strconv"
 "strings"
 "sync"
 "time"
)

type Sample struct {
 Name string `json:"name"`
 Hostname string `json:"hostname"`
 OS string `json:"os"`
 Arch string `json:"arch"`
 CPUCores int `json:"cpu_cores,omitempty"`
 MemoryTotal uint64 `json:"memory_total,omitempty"`
 MemoryUsed uint64 `json:"memory_used,omitempty"`
 DiskTotal uint64 `json:"disk_total,omitempty"`
 DiskUsed uint64 `json:"disk_used,omitempty"`
 CPU float64 `json:"cpu"`
 Memory float64 `json:"memory"`
 Disk float64 `json:"disk"`
 RxBytes uint64 `json:"rx_bytes"`
 TxBytes uint64 `json:"tx_bytes"`
 Uptime uint64 `json:"uptime"`
 Timestamp time.Time `json:"timestamp"`
}
type Node struct { Sample; LastSeen time.Time `json:"last_seen"`; RxSpeed float64 `json:"rx_speed"`; TxSpeed float64 `json:"tx_speed"` }
type Store struct { sync.RWMutex; Nodes map[string]Node `json:"nodes"`; file string; db *sql.DB }
func (s *Store) save() error {
 s.RLock(); b,e:=json.MarshalIndent(s.Nodes,"","  ");s.RUnlock();if e!=nil{return e}
 if e=os.MkdirAll(filepath.Dir(s.file),0700);e!=nil{return e}
 f,e:=os.CreateTemp(filepath.Dir(s.file),".monitor-*");if e!=nil{return e}
 defer os.Remove(f.Name())
 if e=f.Chmod(0600);e!=nil{f.Close();return e}
 if _,e=f.Write(b);e!=nil{f.Close();return e}
 if e=f.Sync();e!=nil{f.Close();return e}
 if e=f.Close();e!=nil{return e}
 return os.Rename(f.Name(),s.file)
}
func env(k, fallback string)string {if x:=os.Getenv(k);x!=""{return x};return fallback}
func equal(a,b string)bool{return subtle.ConstantTimeCompare([]byte(a),[]byte(b))==1}
func main(){
 mode:=flag.String("mode","server","server or agent")
 listen:=flag.String("listen",env("MONITOR_LISTEN","127.0.0.1:8090"),"HTTP listen address")
 server:=flag.String("server","","agent server URL, e.g. https://monitor.example.com")
 name:=flag.String("name","","agent display name")
 token:=flag.String("token","","shared ingestion token (or MONITOR_AGENT_TOKEN)")
 interval:=flag.Duration("interval",5*time.Second,"agent report interval")
 file:=flag.String("data","/var/lib/monitor/nodes.json","legacy snapshot file")
 dbpath:=flag.String("db","/var/lib/monitor/monitor.db","SQLite database file")
 flag.Parse()
 key:=*token;if key==""{key=os.Getenv("MONITOR_AGENT_TOKEN")}
 if key==""{log.Fatal("set -token or MONITOR_AGENT_TOKEN")}
 if *mode=="agent" {runAgent(*server,*name,key,*interval);return}
 if *mode!="server"{log.Fatal("unknown mode")}
 admin:=os.Getenv("MONITOR_ADMIN_TOKEN")
 s:=&Store{Nodes:map[string]Node{},file:*file}
 db,e:=openDB(*dbpath);if e!=nil{log.Fatal(e)};defer db.Close();s.db=db
 if e:=alertSchema(db);e!=nil{log.Fatal(e)}
 if b,e:=os.ReadFile(*file);e==nil{if e=json.Unmarshal(b,&s.Nodes);e!=nil{log.Printf("invalid data file: %v",e)}}
 if e:=restoreNodes(s);e!=nil{log.Printf("restore nodes: %v",e)}
 auth,err:=newAuth(db,admin);if err!=nil{log.Fatal(err)}
 mux:=http.NewServeMux()
 mux.Handle("/static/",staticHandler())
 mux.HandleFunc("/healthz",func(w http.ResponseWriter,r *http.Request){w.Write([]byte("ok"))})
 mux.HandleFunc("/api/v1/ingest",func(w http.ResponseWriter,r *http.Request){
  if r.Method!="POST"{http.Error(w,"method",405);return}
  if !strings.HasPrefix(r.Header.Get("Authorization"),"Bearer "){http.Error(w,"unauthorized",401);return}
  r.Body=http.MaxBytesReader(w,r.Body,16*1024)
  defer r.Body.Close()
  var sample Sample
  if e:=json.NewDecoder(r.Body).Decode(&sample);e!=nil{http.Error(w,"invalid JSON",400);return}
  if len(sample.Name)<1||len(sample.Name)>100||sample.CPU<0||sample.CPU>100||sample.Memory<0||sample.Memory>100||sample.Disk<0||sample.Disk>100||sample.CPUCores<0||sample.CPUCores>4096||sample.MemoryUsed>sample.MemoryTotal||sample.DiskUsed>sample.DiskTotal{http.Error(w,"invalid sample",400);return}
  sample.Timestamp=time.Now().UTC()
  if !validNodeToken(r.Context(),s.db,sample.Name,strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "),key){http.Error(w,"unauthorized",401);return}
  if e:=recordSample(s.db,sample);e!=nil{log.Printf("db: %v",e);http.Error(w,"db write failed",500);return}
  s.Lock();previous,exists:=s.Nodes[sample.Name];node:=Node{Sample:sample,LastSeen:sample.Timestamp};if exists {dt:=sample.Timestamp.Sub(previous.LastSeen).Seconds();if dt>0&&dt<120&&sample.Uptime>=previous.Uptime {if sample.RxBytes>=previous.RxBytes{node.RxSpeed=float64(sample.RxBytes-previous.RxBytes)/dt};if sample.TxBytes>=previous.TxBytes{node.TxSpeed=float64(sample.TxBytes-previous.TxBytes)/dt}}};s.Nodes[sample.Name]=node;s.Unlock()
  w.WriteHeader(204)
 })
 authorized:=auth.require
 auth.routes(mux)
 mux.HandleFunc("/api/v1/tokens",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return};if r.Method!="POST"{http.Error(w,"method",405);return}
  if r.Header.Get("Origin")!="" && r.Header.Get("Origin")!="https://"+r.Host && r.Header.Get("Origin")!="http://"+r.Host{http.Error(w,"origin",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,4096);var req struct{Name string `json:"name"`};if json.NewDecoder(r.Body).Decode(&req)!=nil||!validNodeName(req.Name){http.Error(w,"invalid name",400);return}
 b:=make([]byte,32);if _,e:=rand.Read(b);e!=nil{http.Error(w,"random failed",500);return};secret:=hex.EncodeToString(b);hash:=sha256.Sum256([]byte(secret))
 if _,e:=s.db.ExecContext(r.Context(),"INSERT INTO agent_tokens(name, hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",req.Name,hex.EncodeToString(hash[:]));e!=nil{http.Error(w,"db failed",500);return}
 s.Lock();if _,exists:=s.Nodes[req.Name];!exists{s.Nodes[req.Name]=Node{Sample:Sample{Name:req.Name}}};s.Unlock()
 w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(map[string]string{"name":req.Name,"token":secret})
 })
 mux.HandleFunc("/api/v1/history",func(w http.ResponseWriter,r *http.Request){
  if !authorized(w,r){return}; name:=r.URL.Query().Get("name");hours,_:=strconv.Atoi(r.URL.Query().Get("hours"));if hours!=24&&hours!=168&&hours!=720{hours=24};if len(name)==0||len(name)>100{http.Error(w,"invalid name",400);return}
  points,e:=queryHistory(r.Context(),s.db,name,hours);if e!=nil{http.Error(w,"database error",500);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(points)
 })
 mux.HandleFunc("/api/v1/node-manage",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return};if r.Method!="POST"{http.Error(w,"method",405);return};if !sameOrigin(r){http.Error(w,"origin",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,4096);var x struct{Action string `json:"action"`;Name string `json:"name"`;NewName string `json:"new_name"`};if json.NewDecoder(r.Body).Decode(&x)!=nil||!validNodeName(x.Name){http.Error(w,"invalid input",400);return}
 s.Lock();defer s.Unlock();
 switch x.Action {
 case "delete":
   if e:=removeNode(r.Context(),s.db,x.Name);e!=nil{http.Error(w,e.Error(),500);return};delete(s.Nodes,x.Name)
 case "revoke":
   if e:=revokeNode(r.Context(),s.db,x.Name);e!=nil{http.Error(w,e.Error(),500);return};delete(s.Nodes,x.Name)
 case "rename":
   if !validNodeName(x.NewName)||x.NewName==x.Name{http.Error(w,"invalid new name",400);return}
   if _,ok:=s.Nodes[x.NewName];ok{http.Error(w,"name exists",409);return}
   if e:=renameNode(r.Context(),s.db,x.Name,x.NewName);e!=nil{http.Error(w,e.Error(),500);return}
   if n,ok:=s.Nodes[x.Name];ok{n.Name=x.NewName;s.Nodes[x.NewName]=n;delete(s.Nodes,x.Name)}
 default:http.Error(w,"invalid action",400);return
 };w.WriteHeader(204)
 })
 mux.HandleFunc("/api/v1/node-metadata",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return}
 if r.Method!="POST"{http.Error(w,"method",405);return}
 if !sameOrigin(r){http.Error(w,"origin",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,8192)
 var in struct{Name string `json:"name"`; DisplayName string `json:"display_name"`; Group string `json:"group"`; Notes string `json:"notes"`}
 if json.NewDecoder(r.Body).Decode(&in)!=nil||!validNodeName(in.Name)||len([]rune(in.DisplayName))>60||len([]rune(in.Group))>40||len([]rune(in.Notes))>500{http.Error(w,"invalid metadata",400);return}
 in.DisplayName=strings.TrimSpace(in.DisplayName);in.Group=strings.TrimSpace(in.Group);in.Notes=strings.TrimSpace(in.Notes)
 s.RLock();_,found:=s.Nodes[in.Name];s.RUnlock();if !found{http.Error(w,"unknown node",404);return}
 _,e:=s.db.ExecContext(r.Context(),"INSERT INTO node_metadata(name,display_name,group_name,notes) VALUES(?,?,?,?) ON CONFLICT(name) DO UPDATE SET display_name=excluded.display_name,group_name=excluded.group_name,notes=excluded.notes",in.Name,in.DisplayName,in.Group,in.Notes)
 if e!=nil{http.Error(w,"database error",500);return};w.WriteHeader(204)
 })
 mux.HandleFunc("/api/v1/traffic",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return};name:=r.URL.Query().Get("name");hours,_:=strconv.Atoi(r.URL.Query().Get("hours"));if hours!=24&&hours!=168&&hours!=720{hours=24};if !validNodeName(name){http.Error(w,"invalid name",400);return}
 result,e:=queryTraffic(r.Context(),s.db,name,hours);if e!=nil{http.Error(w,"database error",500);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(result)
 })
 mux.HandleFunc("/api/v1/alerts",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return}
 w.Header().Set("Cache-Control","no-store");w.Header().Set("Content-Type","application/json")
 if r.Method=="GET"{items,e:=readAlerts(r.Context(),s.db);if e!=nil{http.Error(w,"database error",500);return};json.NewEncoder(w).Encode(items);return}
 if r.Method!="POST"||!sameOrigin(r){http.Error(w,"invalid method or origin",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,4096)
 var cfg AlertSettings
 if json.NewDecoder(r.Body).Decode(&cfg)!=nil||cfg.OfflineSeconds<30||cfg.OfflineSeconds>3600||cfg.CPUThreshold<1||cfg.CPUThreshold>100||cfg.MemoryThreshold<1||cfg.MemoryThreshold>100||cfg.DiskThreshold<1||cfg.DiskThreshold>100||cfg.DurationSeconds<30||cfg.DurationSeconds>3600||len(cfg.Webhook)>512 {http.Error(w,"invalid settings",400);return}
 if cfg.Webhook!=""&&!validWebhook(cfg.Webhook){http.Error(w,"webhook must be HTTPS without credentials",400);return}
 if e:=writeAlertSettings(s.db,cfg);e!=nil{http.Error(w,"database error",500);return};w.WriteHeader(204)
 })
 mux.HandleFunc("/api/v1/nodes",func(w http.ResponseWriter,r *http.Request){
  if !authorized(w,r){return}
  meta:=map[string][3]string{};rows,e:=s.db.QueryContext(r.Context(),"SELECT name,display_name,group_name,notes FROM node_metadata");if e!=nil{http.Error(w,"database error",500);return};for rows.Next(){var name,display,group,notes string;if rows.Scan(&name,&display,&group,&notes)==nil{meta[name]=[3]string{display,group,notes}}};rows.Close()
  s.RLock();out:=make([]map[string]any,0,len(s.Nodes))
  for _,n:=range s.Nodes{out=append(out,map[string]any{"name":n.Name,"display_name":meta[n.Name][0],"group":meta[n.Name][1],"notes":meta[n.Name][2],"hostname":n.Hostname,"os":n.OS,"arch":n.Arch,"cpu":n.CPU,"cpu_cores":n.CPUCores,"memory_total":n.MemoryTotal,"memory_used":n.MemoryUsed,"disk_total":n.DiskTotal,"disk_used":n.DiskUsed,"memory":n.Memory,"disk":n.Disk,"rx_bytes":n.RxBytes,"tx_bytes":n.TxBytes,"uptime":n.Uptime,"rx_speed":n.RxSpeed,"tx_speed":n.TxSpeed,"last_seen":n.LastSeen,"online":time.Since(n.LastSeen)<30*time.Second})};s.RUnlock()
  w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(out)
 })
 mux.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/"{http.NotFound(w,r);return}
  if !authorized(w,r){http.Redirect(w,r,"/login",303);return}
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.Header().Set("Content-Security-Policy","default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
  dashboard.Execute(w,nil)
 })
 go startAlertLoop(s)
 go func(){t:=time.NewTicker(30*time.Second);defer t.Stop();for range t.C{if e:=s.save();e!=nil{log.Printf("save: %v",e)};if _,e:=s.db.Exec("DELETE FROM samples WHERE ts < ?",time.Now().Add(-30*24*time.Hour).Unix());e!=nil{log.Printf("retention: %v",e)}}}()
 log.Printf("Monitor %s server listening on %s (%s)",runtime.Version(),*listen,*file)
 srv:=&http.Server{Addr:*listen,Handler:mux,ReadHeaderTimeout:5*time.Second,ReadTimeout:10*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second}
 log.Fatal(srv.ListenAndServe())
}
var dashboard=template.Must(template.New("dash").Parse(dashboardHTML))
var _=fmt.Sprintf
var _=strconv.Itoa
