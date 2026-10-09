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
 CPU float64 `json:"cpu"`
 Memory float64 `json:"memory"`
 Disk float64 `json:"disk"`
 RxBytes uint64 `json:"rx_bytes"`
 TxBytes uint64 `json:"tx_bytes"`
 Uptime uint64 `json:"uptime"`
 Timestamp time.Time `json:"timestamp"`
}
type Node struct { Sample; LastSeen time.Time `json:"last_seen"` }
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
 listen:=flag.String("listen","127.0.0.1:8090","HTTP listen address")
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
 admin:=os.Getenv("MONITOR_ADMIN_TOKEN");if len(admin)<24{log.Fatal("set MONITOR_ADMIN_TOKEN (at least 24 characters)")}
 s:=&Store{Nodes:map[string]Node{},file:*file}
 db,e:=openDB(*dbpath);if e!=nil{log.Fatal(e)};defer db.Close();s.db=db
 if b,e:=os.ReadFile(*file);e==nil{if e=json.Unmarshal(b,&s.Nodes);e!=nil{log.Printf("invalid data file: %v",e)}}
 if e:=restoreNodes(s);e!=nil{log.Printf("restore nodes: %v",e)}
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
  if len(sample.Name)<1||len(sample.Name)>100||sample.CPU<0||sample.CPU>100||sample.Memory<0||sample.Memory>100||sample.Disk<0||sample.Disk>100{http.Error(w,"invalid sample",400);return}
  sample.Timestamp=time.Now().UTC()
  if !validNodeToken(r.Context(),s.db,sample.Name,strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "),key){http.Error(w,"unauthorized",401);return}
  if e:=recordSample(s.db,sample);e!=nil{log.Printf("db: %v",e);http.Error(w,"db write failed",500);return}
  s.Lock();s.Nodes[sample.Name]=Node{Sample:sample,LastSeen:sample.Timestamp};s.Unlock()
  w.WriteHeader(204)
 })
 authorized:=func(w http.ResponseWriter,r *http.Request)bool{
  c,e:=r.Cookie("monitor_session")
  if e!=nil||!equal(c.Value,admin){http.Error(w,"unauthorized",401);return false};return true
 }
 mux.HandleFunc("/api/v1/tokens",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return};if r.Method!="POST"{http.Error(w,"method",405);return}
  if r.Header.Get("Origin")!="" && r.Header.Get("Origin")!="https://"+r.Host && r.Header.Get("Origin")!="http://"+r.Host{http.Error(w,"origin",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,4096);var req struct{Name string `json:"name"`};if json.NewDecoder(r.Body).Decode(&req)!=nil||!validNodeName(req.Name){http.Error(w,"invalid name",400);return}
 b:=make([]byte,32);if _,e:=rand.Read(b);e!=nil{http.Error(w,"random failed",500);return};secret:=hex.EncodeToString(b);hash:=sha256.Sum256([]byte(secret))
 if _,e:=s.db.ExecContext(r.Context(),"INSERT INTO agent_tokens(name, hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",req.Name,hex.EncodeToString(hash[:]));e!=nil{http.Error(w,"db failed",500);return}
 s.Lock();if _,exists:=s.Nodes[req.Name];!exists{s.Nodes[req.Name]=Node{Sample:Sample{Name:req.Name}}};s.Unlock()
 w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(map[string]string{"name":req.Name,"token":secret})
 })
 mux.HandleFunc("/login",func(w http.ResponseWriter,r *http.Request){
  if r.Method=="GET"{w.Header().Set("Content-Type","text/html; charset=utf-8");w.Write([]byte(loginHTML));return}
  if r.Method!="POST"{http.Error(w,"method",405);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096)
  if e:=r.ParseForm();e!=nil||!equal(r.FormValue("token"),admin){http.Error(w,"invalid login",401);return}
  http.SetCookie(w,&http.Cookie{Name:"monitor_session",Value:admin,Path:"/",HttpOnly:true,Secure:r.TLS!=nil||r.Header.Get("X-Forwarded-Proto")=="https",SameSite:http.SameSiteStrictMode})
  http.Redirect(w,r,"/",303)
 })
 mux.HandleFunc("/logout",func(w http.ResponseWriter,r *http.Request){
  if r.Method!="POST"||!authorized(w,r){return}
  http.SetCookie(w,&http.Cookie{Name:"monitor_session",Path:"/",MaxAge:-1,HttpOnly:true,SameSite:http.SameSiteStrictMode})
  http.Redirect(w,r,"/login",303)
 })
 mux.HandleFunc("/api/v1/history",func(w http.ResponseWriter,r *http.Request){
  if !authorized(w,r){return}; name:=r.URL.Query().Get("name");hours,_:=strconv.Atoi(r.URL.Query().Get("hours"));if hours!=24&&hours!=168&&hours!=720{hours=24};if len(name)==0||len(name)>100{http.Error(w,"invalid name",400);return}
  points,e:=queryHistory(r.Context(),s.db,name,hours);if e!=nil{http.Error(w,"database error",500);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(points)
 })
 mux.HandleFunc("/api/v1/traffic",func(w http.ResponseWriter,r *http.Request){
 if !authorized(w,r){return};name:=r.URL.Query().Get("name");hours,_:=strconv.Atoi(r.URL.Query().Get("hours"));if hours!=24&&hours!=168&&hours!=720{hours=24};if !validNodeName(name){http.Error(w,"invalid name",400);return}
 result,e:=queryTraffic(r.Context(),s.db,name,hours);if e!=nil{http.Error(w,"database error",500);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(result)
 })
 mux.HandleFunc("/api/v1/nodes",func(w http.ResponseWriter,r *http.Request){
  if !authorized(w,r){return}
  s.RLock();out:=make([]map[string]any,0,len(s.Nodes))
  for _,n:=range s.Nodes{out=append(out,map[string]any{"name":n.Name,"hostname":n.Hostname,"os":n.OS,"arch":n.Arch,"cpu":n.CPU,"memory":n.Memory,"disk":n.Disk,"rx_bytes":n.RxBytes,"tx_bytes":n.TxBytes,"uptime":n.Uptime,"last_seen":n.LastSeen,"online":time.Since(n.LastSeen)<30*time.Second})};s.RUnlock()
  w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");json.NewEncoder(w).Encode(out)
 })
 mux.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/"{http.NotFound(w,r);return}
  if !authorized(w,r){return}
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.Header().Set("Content-Security-Policy","default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
  dashboard.Execute(w,nil)
 })
 go func(){t:=time.NewTicker(30*time.Second);defer t.Stop();for range t.C{if e:=s.save();e!=nil{log.Printf("save: %v",e)};if _,e:=s.db.Exec("DELETE FROM samples WHERE ts < ?",time.Now().Add(-30*24*time.Hour).Unix());e!=nil{log.Printf("retention: %v",e)}}}()
 log.Printf("Monitor %s server listening on %s (%s)",runtime.Version(),*listen,*file)
 srv:=&http.Server{Addr:*listen,Handler:mux,ReadHeaderTimeout:5*time.Second,ReadTimeout:10*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second}
 log.Fatal(srv.ListenAndServe())
}
var dashboard=template.Must(template.New("dash").Parse(dashboardHTML))
var _=fmt.Sprintf
var _=strconv.Itoa
