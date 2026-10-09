package main

import (
 "crypto/rand"
 "crypto/sha256"
 "crypto/subtle"
 "database/sql"
 "encoding/hex"
 "errors"
 "fmt"
 "html/template"
 "net"
 "net/http"
 "strings"
 "sync"
 "time"
 "golang.org/x/crypto/bcrypt"
)
type attempt struct {failures int; blocked time.Time}
type Auth struct {db *sql.DB; legacy string; mu sync.Mutex; tries map[string]attempt}
func newAuth(db *sql.DB,legacy string)(*Auth,error){
 for _,q:=range []string{
  "CREATE TABLE IF NOT EXISTS admin_user (id INTEGER PRIMARY KEY CHECK (id=1),username TEXT NOT NULL,password_hash TEXT NOT NULL)",
  "CREATE TABLE IF NOT EXISTS admin_sessions (token_hash TEXT PRIMARY KEY,expires INTEGER NOT NULL)",
 }{if _,e:=db.Exec(q);e!=nil{return nil,e}}
 return &Auth{db:db,legacy:legacy,tries:map[string]attempt{}},nil
}
func (a *Auth) configured()bool{var n int;return a.db.QueryRow("SELECT COUNT(*) FROM admin_user").Scan(&n)==nil&&n>0}
func clientIP(r *http.Request)string{h,_,e:=net.SplitHostPort(r.RemoteAddr);if e!=nil{return r.RemoteAddr};return h}
func (a *Auth) limited(ip string)bool{a.mu.Lock();defer a.mu.Unlock();v:=a.tries[ip];return time.Now().Before(v.blocked)}
func (a *Auth) failed(ip string){a.mu.Lock();defer a.mu.Unlock();v:=a.tries[ip];v.failures++;if v.failures>=5{v.blocked=time.Now().Add(15*time.Minute);v.failures=0};a.tries[ip]=v}
func (a *Auth) clear(ip string){a.mu.Lock();delete(a.tries,ip);a.mu.Unlock()}
func randomSecret()(string,error){b:=make([]byte,32);_,e:=rand.Read(b);return hex.EncodeToString(b),e}
func digest(s string)string{h:=sha256.Sum256([]byte(s));return hex.EncodeToString(h[:])}
func secure(r *http.Request)bool{return r.TLS!=nil || r.Header.Get("X-Forwarded-Proto")=="https"}
func (a *Auth) issue(w http.ResponseWriter,r *http.Request)error{
 secret,e:=randomSecret();if e!=nil{return e}
 _,e=a.db.Exec("INSERT INTO admin_sessions(token_hash,expires) VALUES(?,?)",digest(secret),time.Now().Add(7*24*time.Hour).Unix());if e!=nil{return e}
 http.SetCookie(w,&http.Cookie{Name:"monitor_session",Value:secret,Path:"/",HttpOnly:true,Secure:secure(r),SameSite:http.SameSiteStrictMode,MaxAge:604800})
 return nil
}
func (a *Auth) session(r *http.Request)bool{
 c,e:=r.Cookie("monitor_session");if e!=nil||len(c.Value)!=64{return false}
 var expires int64
 if a.db.QueryRow("SELECT expires FROM admin_sessions WHERE token_hash=?",digest(c.Value)).Scan(&expires)!=nil{return false}
 return time.Now().Unix()<expires
}
func (a *Auth) require(w http.ResponseWriter,r *http.Request)bool{
 if !a.session(r){http.Error(w,"unauthorized",401);return false};return true
}
const authHTML=`<!doctype html><html lang="zh"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Monitor 账户</title><style>body{font:16px system-ui;background:#0b1220;color:#e8f0ff;min-height:85vh;display:grid;place-items:center}main{width:min(390px,86vw);background:#17243a;padding:26px;border-radius:15px}input,button{width:100%;padding:12px;margin:7px 0;box-sizing:border-box;border-radius:7px;border:1px solid #56647a;background:#0b1220;color:white}button{background:#2679d9;border:none}p{color:#acb9ca;font-size:13px}</style></head><body><main><h2>{{.Title}}</h2><form method="post"><input name="username" placeholder="管理员账号" required autocomplete="username" maxlength="64"><input name="password" type="password" placeholder="密码（至少 12 位）" required autocomplete="{{.Autocomplete}}">{{if .Setup}}<input name="legacy" type="password" placeholder="原管理员令牌（迁移验证）" required>{{else}}<input name="old_password" type="hidden">{{end}}<button>{{.Action}}</button></form><p>{{.Hint}}</p></main></body></html>`
func authPage(w http.ResponseWriter,setup bool){
 w.Header().Set("Content-Type","text/html; charset=utf-8");w.Header().Set("Cache-Control","no-store")
 h:=template.Must(template.New("auth").Parse(authHTML))
 v:=map[string]any{"Title":"Monitor 登录","Action":"登录","Hint":"请使用管理员账号和密码","Autocomplete":"current-password","Setup":false}
 if setup{v["Title"]="初始化管理员";v["Action"]="创建账号";v["Hint"]="升级用户需要原管理员令牌；请在 HTTPS 页面完成初始化。";v["Autocomplete"]="new-password";v["Setup"]=true}
 h.Execute(w,v)
}
func validUsername(s string)bool{if len(s)<3||len(s)>64{return false};for _,c:=range s{if !(c>='a'&&c<='z'||c>='A'&&c<='Z'||c>='0'&&c<='9'||c=='_'||c=='-'){return false}};return true}
func (a *Auth) routes(mux *http.ServeMux){
 mux.HandleFunc("/login",func(w http.ResponseWriter,r *http.Request){
  if !a.configured(){http.Redirect(w,r,"/setup",303);return}
  if r.Method=="GET"{authPage(w,false);return};if r.Method!="POST"{http.Error(w,"method",405);return}
  if !sameOrigin(r){http.Error(w,"origin",403);return}
  ip:=clientIP(r);if a.limited(ip){http.Error(w,"too many attempts",429);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096);if r.ParseForm()!=nil{http.Error(w,"invalid",400);return}
  var username,hash string
  e:=a.db.QueryRow("SELECT username,password_hash FROM admin_user WHERE id=1").Scan(&username,&hash)
  if e!=nil||subtle.ConstantTimeCompare([]byte(username),[]byte(r.FormValue("username")))!=1||bcrypt.CompareHashAndPassword([]byte(hash),[]byte(r.FormValue("password")))!=nil{a.failed(ip);http.Error(w,"账号或密码错误",401);return}
  a.clear(ip);if a.issue(w,r)!=nil{http.Error(w,"session error",500);return};http.Redirect(w,r,"/",303)
 })
 mux.HandleFunc("/setup",func(w http.ResponseWriter,r *http.Request){
  if a.configured(){http.Redirect(w,r,"/login",303);return}
  if r.Method=="GET"{authPage(w,true);return};if r.Method!="POST"{http.Error(w,"method",405);return}
  if !sameOrigin(r){http.Error(w,"origin",403);return}
  ip:=clientIP(r);if a.limited(ip){http.Error(w,"too many attempts",429);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096);if r.ParseForm()!=nil{http.Error(w,"invalid",400);return}
  name,pass:=r.FormValue("username"),r.FormValue("password")
  if !validUsername(name)||len(pass)<12||len(pass)>72{http.Error(w,"账号或密码不符合要求",400);return}
  if len(a.legacy)<24||subtle.ConstantTimeCompare([]byte(a.legacy),[]byte(r.FormValue("legacy")))!=1{a.failed(ip);http.Error(w,"初始化令牌错误",401);return}
  hash,e:=bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost);if e!=nil{http.Error(w,"hash error",500);return}
  _,e=a.db.Exec("INSERT INTO admin_user(id,username,password_hash) VALUES(1,?,?)",name,string(hash));if e!=nil{http.Error(w,"already initialized",409);return}
  a.clear(ip);if a.issue(w,r)!=nil{http.Error(w,"session error",500);return};http.Redirect(w,r,"/",303)
 })
 mux.HandleFunc("/logout",func(w http.ResponseWriter,r *http.Request){
  if r.Method!="POST"||!sameOrigin(r)||!a.require(w,r){return}
  if c,e:=r.Cookie("monitor_session");e==nil{a.db.Exec("DELETE FROM admin_sessions WHERE token_hash=?",digest(c.Value))}
  http.SetCookie(w,&http.Cookie{Name:"monitor_session",Path:"/",MaxAge:-1,HttpOnly:true,Secure:secure(r),SameSite:http.SameSiteStrictMode})
  http.Redirect(w,r,"/login",303)
 })
 mux.HandleFunc("/password",func(w http.ResponseWriter,r *http.Request){
  if !a.require(w,r){return};if r.Method=="GET"{w.Header().Set("Content-Type","text/html; charset=utf-8");fmt.Fprint(w,`<!doctype html><html lang="zh"><meta charset="utf-8"><title>修改密码</title><form method="post"><input type="password" name="old" placeholder="原密码" required><input type="password" name="new" placeholder="新密码（至少12位）" required><button>修改密码</button></form></html>`);return}
  if r.Method!="POST"||!sameOrigin(r){http.Error(w,"invalid",403);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096);if r.ParseForm()!=nil{http.Error(w,"invalid",400);return}
  var current string;if e:=a.db.QueryRow("SELECT password_hash FROM admin_user WHERE id=1").Scan(&current);e!=nil{http.Error(w,"db error",500);return}
  newPass:=r.FormValue("new");if len(newPass)<12||len(newPass)>72||bcrypt.CompareHashAndPassword([]byte(current),[]byte(r.FormValue("old")))!=nil{http.Error(w,"旧密码错误或新密码长度无效",400);return}
  hash,e:=bcrypt.GenerateFromPassword([]byte(newPass),bcrypt.DefaultCost);if e!=nil{http.Error(w,"hash error",500);return}
  tx,e:=a.db.Begin();if e!=nil{http.Error(w,"db error",500);return};defer tx.Rollback()
  if _,e=tx.Exec("UPDATE admin_user SET password_hash=? WHERE id=1",string(hash));e!=nil{http.Error(w,"db error",500);return}
  if _,e=tx.Exec("DELETE FROM admin_sessions");e!=nil{http.Error(w,"db error",500);return}
  if e=tx.Commit();e!=nil{http.Error(w,"db error",500);return}
  http.SetCookie(w,&http.Cookie{Name:"monitor_session",Path:"/",MaxAge:-1,HttpOnly:true,Secure:secure(r),SameSite:http.SameSiteStrictMode})
  http.Redirect(w,r,"/login",303)
 })
}
var _=errors.New
