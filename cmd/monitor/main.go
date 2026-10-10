package main

import (
 "crypto/subtle"
 
 "encoding/json"
 "database/sql"
 "crypto/sha256"
 "encoding/hex"
 "crypto/rand"
 _ "modernc.org/sqlite"
 _ "time/tzdata"
 "flag"
 "fmt"
 "html/template"
 "log"
 "net"
 "net/http"
 "os"
 "path/filepath"
 "runtime"
 "strconv"
 "strings"
 "sync"
 "time"
)


func validIPFamily(raw string,family int)bool{
 if !validPublicIP(raw){return false}
 parsed:=net.ParseIP(raw)
 if parsed==nil{return false}
 if family==4{return parsed.To4()!=nil}
 if family==6{return parsed.To4()==nil && parsed.To16()!=nil}
 return false
}

const monitorVersion="v0.9.7"
type Sample struct {
 Name string `json:"name"`
 Hostname string `json:"hostname"`
 OS string `json:"os"`
 Arch string `json:"arch"`
 AgentVersion string `json:"agent_version,omitempty"`
 PublicIP string `json:"public_ip,omitempty"`
 AutoLocation string `json:"auto_location,omitempty"`
 CountryCode string `json:"country_code,omitempty"`
 PublicIPv4 string `json:"public_ipv4,omitempty"`
 PublicIPv6 string `json:"public_ipv6,omitempty"`
 CPUCores int `json:"cpu_cores,omitempty"`
 CPUModel string `json:"cpu_model,omitempty"`
 Load1 float64 `json:"load1,omitempty"`
 Load5 float64 `json:"load5,omitempty"`
 Load15 float64 `json:"load15,omitempty"`
 SwapTotal uint64 `json:"swap_total,omitempty"`
 SwapUsed uint64 `json:"swap_used,omitempty"`
 DiskReadBytes uint64 `json:"disk_read_bytes,omitempty"`
 DiskWriteBytes uint64 `json:"disk_write_bytes,omitempty"`
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
type Node struct { Sample; LastSeen time.Time `json:"last_seen"`; RxSpeed float64 `json:"rx_speed"`; TxSpeed float64 `json:"tx_speed"`;DiskReadSpeed float64 `json:"disk_read_speed"`;DiskWriteSpeed float64 `json:"disk_write_speed"` }
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
  if len(sample.Name)<1||len(sample.Name)>100||sample.CPU<0||sample.CPU>100||sample.Memory<0||sample.Memory>100||sample.Disk<0||sample.Disk>100||sample.CPUCores<0||sample.CPUCores>4096||sample.MemoryUsed>sample.MemoryTotal||sample.DiskUsed>sample.DiskTotal||sample.SwapUsed>sample.SwapTotal||len(sample.CPUModel)>256||len(sample.AgentVersion)>32||len([]rune(sample.AutoLocation))>140||len(sample.PublicIP)>45||len(sample.CountryCode)>2||len(sample.PublicIPv4)>15||len(sample.PublicIPv6)>45||(sample.PublicIPv4!=""&&!validIPFamily(sample.PublicIPv4,4))||(sample.PublicIPv6!=""&&!validIPFamily(sample.PublicIPv6,6)){http.Error(w,"invalid sample",400);return}
  sample.Timestamp=time.Now().UTC()
  if !validNodeToken(r.Context(),s.db,sample.Name,strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "),key){http.Error(w,"unauthorized",401);return}
  // Only store IP geolocation supplied by the authenticated Agent. A
  // Cloudflare/Nginx reverse proxy would otherwise appear as the node.
  if sample.PublicIP!=""&&sample.AutoLocation!=""&&validPublicIP(sample.PublicIP){
   if _,e:=s.db.ExecContext(r.Context(),"INSERT INTO node_geo(name,public_ip,auto_location,country_code,updated_at) VALUES(?,?,?,?,?) ON CONFLICT(name) DO UPDATE SET public_ip=excluded.public_ip,auto_location=excluded.auto_location,country_code=excluded.country_code,updated_at=excluded.updated_at WHERE public_ip<>excluded.public_ip OR auto_location<>excluded.auto_location OR country_code<>excluded.country_code",sample.Name,sample.PublicIP,sample.AutoLocation,sample.CountryCode,sample.Timestamp.Unix());e!=nil{log.Printf("geo save: %v",e)}
  }

  if e:=recordTrafficIncrement(s.db,sample);e!=nil{log.Printf("traffic: %v",e)}
  if e:=recordSample(s.db,sample);e!=nil{log.Printf("db: %v",e);http.Error(w,"db write failed",500);return}
  s.Lock();previous,exists:=s.Nodes[sample.Name];node:=Node{Sample:sample,LastSeen:sample.Timestamp};if exists {dt:=sample.Timestamp.Sub(previous.LastSeen).Seconds();if dt>0&&dt<120&&sample.Uptime>=previous.Uptime {if sample.RxBytes>=previous.RxBytes{node.RxSpeed=float64(sample.RxBytes-previous.RxBytes)/dt};if sample.TxBytes>=previous.TxBytes{node.TxSpeed=float64(sample.TxBytes-previous.TxBytes)/dt};if sample.DiskReadBytes>=previous.DiskReadBytes{node.DiskReadSpeed=float64(sample.DiskReadBytes-previous.DiskReadBytes)/dt};if sample.DiskWriteBytes>=previous.DiskWriteBytes{node.DiskWriteSpeed=float64(sample.DiskWriteBytes-previous.DiskWriteBytes)/dt}}};s.Nodes[sample.Name]=node;s.Unlock()
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
 r.Body=http.MaxBytesReader(w,r.Body,8192);defer r.Body.Close()
 var in struct{
  Name string `json:"name"`
  DisplayName string `json:"display_name"`
  Group string `json:"group"`
  Location string `json:"location"`
  Notes string `json:"notes"`
  Profile *NodeProfile `json:"profile"`
  QuotaGB *float64 `json:"quota_gb"`
  ExpiresOn *string `json:"expires_on"`
 }
 if json.NewDecoder(r.Body).Decode(&in)!=nil||!validNodeName(in.Name)||len([]rune(in.DisplayName))>60||len([]rune(in.Group))>40||len([]rune(in.Location))>80||len([]rune(in.Notes))>500{
  http.Error(w,"invalid metadata",400);return
 }
 in.DisplayName=strings.TrimSpace(in.DisplayName);in.Group=strings.TrimSpace(in.Group);in.Location=strings.TrimSpace(in.Location);in.Notes=strings.TrimSpace(in.Notes)
 if in.Profile!=nil{
  in.Profile.Provider=strings.TrimSpace(in.Profile.Provider)
  in.Profile.CountryCode=strings.ToUpper(strings.TrimSpace(in.Profile.CountryCode))
  if !validNodeProfile(*in.Profile){http.Error(w,"invalid server plan",400);return}
 }
 if (in.QuotaGB==nil)!=(in.ExpiresOn==nil){http.Error(w,"quota and expiry must be supplied together",400);return}
 s.RLock();_,found:=s.Nodes[in.Name];s.RUnlock();if !found{http.Error(w,"unknown node",404);return}
 tx,e:=s.db.BeginTx(r.Context(),nil);if e!=nil{http.Error(w,"database error",500);return};defer tx.Rollback()
 if _,e=tx.ExecContext(r.Context(),"INSERT INTO node_metadata(name,display_name,group_name,location,notes) VALUES(?,?,?,?,?) ON CONFLICT(name) DO UPDATE SET display_name=excluded.display_name,group_name=excluded.group_name,location=excluded.location,notes=excluded.notes",in.Name,in.DisplayName,in.Group,in.Location,in.Notes);e!=nil{http.Error(w,"database error",500);return}
 if in.Profile!=nil{
  p:=in.Profile
  _,e=tx.ExecContext(r.Context(),"INSERT INTO node_profile(name,provider,country_code,price_value,price_currency,billing_cycle,period_start,port_mbps,has_ipv4,has_ipv6) VALUES(?,?,?,?,?,?,?,?,?,?) ON CONFLICT(name) DO UPDATE SET provider=excluded.provider,country_code=excluded.country_code,price_value=excluded.price_value,price_currency=excluded.price_currency,billing_cycle=excluded.billing_cycle,period_start=excluded.period_start,port_mbps=excluded.port_mbps,has_ipv4=excluded.has_ipv4,has_ipv6=excluded.has_ipv6",in.Name,p.Provider,p.CountryCode,p.PriceValue,p.PriceCurrency,p.BillingCycle,p.PeriodStart,p.PortMbps,p.HasIPv4,p.HasIPv6)
  if e!=nil{http.Error(w,"database error",500);return}
 }
 if in.QuotaGB!=nil{
  var zone string
  e=tx.QueryRowContext(r.Context(),"SELECT timezone FROM node_limits WHERE name=?",in.Name).Scan(&zone)
  if e==sql.ErrNoRows{zone="UTC"}else if e!=nil{http.Error(w,"database error",500);return}
  limits:=NodeLimits{Name:in.Name,Timezone:zone,QuotaGB:*in.QuotaGB,ExpiresOn:*in.ExpiresOn}
  if !validNodeLimits(limits){http.Error(w,"invalid quota or expiry",400);return}
  _,e=tx.ExecContext(r.Context(),"INSERT INTO node_limits(name,timezone,quota_gb,expires_on) VALUES(?,?,?,?) ON CONFLICT(name) DO UPDATE SET quota_gb=excluded.quota_gb,expires_on=excluded.expires_on",in.Name,zone,*in.QuotaGB,*in.ExpiresOn)
  if e!=nil{http.Error(w,"database error",500);return}
 }
 if e=tx.Commit();e!=nil{http.Error(w,"database error",500);return}
 w.WriteHeader(http.StatusNoContent)
 })
 mux.HandleFunc("/api/v1/node-limits",func(w http.ResponseWriter,r *http.Request){
  if !authorized(w,r){return}
  if r.Method=="GET"{
   settings,e:=readNodeLimits(r.Context(),s.db);if e!=nil{http.Error(w,"database error",500);return}
   w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(settings);return
  }
  if r.Method!="POST"||!sameOrigin(r){http.Error(w,"invalid method or origin",403);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096);defer r.Body.Close()
  var v NodeLimits
  if json.NewDecoder(r.Body).Decode(&v)!=nil||!validNodeLimits(v){http.Error(w,"invalid timezone, quota or expiry",400);return}
  s.RLock();_,ok:=s.Nodes[v.Name];s.RUnlock();if !ok{http.Error(w,"unknown node",404);return}
  if e:=setNodeLimits(r.Context(),s.db,v);e!=nil{http.Error(w,"database error",500);return}
  w.WriteHeader(http.StatusNoContent)
 })
 mux.HandleFunc("/api/v1/traffic-summary",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return}
 sums,e:=periodTraffic(r.Context(),s.db);if e!=nil{http.Error(w,"database error",500);return}
 w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(sums)
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
  geo:=map[string][3]string{};gr,e:=s.db.QueryContext(r.Context(),"SELECT name,public_ip,auto_location,country_code FROM node_geo");if e!=nil{http.Error(w,"database error",500);return};for gr.Next(){var name,ip,location,country string;if gr.Scan(&name,&ip,&location,&country)==nil{geo[name]=[3]string{ip,location,country}}};if e=gr.Err();e!=nil{gr.Close();http.Error(w,"database error",500);return};gr.Close()
  meta:=map[string][4]string{};rows,e:=s.db.QueryContext(r.Context(),"SELECT name,display_name,group_name,location,notes FROM node_metadata");if e!=nil{http.Error(w,"database error",500);return};for rows.Next(){var name,display,group,location,notes string;if rows.Scan(&name,&display,&group,&location,&notes)==nil{meta[name]=[4]string{display,group,location,notes}}};if e=rows.Err();e!=nil{rows.Close();http.Error(w,"database error",500);return};rows.Close()
  profiles,e:=readNodeProfiles(r.Context(),s.db);if e!=nil{http.Error(w,"database error",500);return}
  s.RLock();out:=make([]map[string]any,0,len(s.Nodes))
  for _,n:=range s.Nodes{p,exists:=profiles[n.Name];if !exists{p=defaultNodeProfile()};country:=p.CountryCode;if country==""{country=geo[n.Name][2]};out=append(out,map[string]any{"profile":p,"country_code":country,"name":n.Name,"display_name":meta[n.Name][0],"group":meta[n.Name][1],"location":func()string{if meta[n.Name][2]!=""{return meta[n.Name][2]};return geo[n.Name][1]}(),"manual_location":meta[n.Name][2],"auto_location":geo[n.Name][1],"public_ip":geo[n.Name][0],"public_ipv4":n.PublicIPv4,"public_ipv6":n.PublicIPv6,"location_source":func()string{if meta[n.Name][2]!=""{return "manual"};if geo[n.Name][1]!=""{return "ip"};return "unknown"}(),"notes":meta[n.Name][3],"hostname":n.Hostname,"os":n.OS,"arch":n.Arch,"cpu":n.CPU,"cpu_cores":n.CPUCores,"cpu_model":n.CPUModel,"agent_version":n.AgentVersion,"load1":n.Load1,"load5":n.Load5,"load15":n.Load15,"swap_total":n.SwapTotal,"swap_used":n.SwapUsed,"disk_read_speed":n.DiskReadSpeed,"disk_write_speed":n.DiskWriteSpeed,"memory_total":n.MemoryTotal,"memory_used":n.MemoryUsed,"disk_total":n.DiskTotal,"disk_used":n.DiskUsed,"memory":n.Memory,"disk":n.Disk,"rx_bytes":n.RxBytes,"tx_bytes":n.TxBytes,"uptime":n.Uptime,"rx_speed":n.RxSpeed,"tx_speed":n.TxSpeed,"last_seen":n.LastSeen,"online":time.Since(n.LastSeen)<30*time.Second})};s.RUnlock()
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
