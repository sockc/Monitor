package main

import (
 "context"
 "crypto/sha256"
 "crypto/subtle"
 "database/sql"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "os"
 "net/http"
 "path/filepath"
 "sort"
 "time"
)
type Point struct{Time int64 `json:"time"`;CPU float64 `json:"cpu"`; Memory float64 `json:"memory"`; Disk float64 `json:"disk"`}
func openDB(path string)(*sql.DB,error){
 if err:=os.MkdirAll(filepath.Dir(path),0700);err!=nil{return nil,err}
 db,err:=sql.Open("sqlite",path);if err!=nil{return nil,err}
 db.SetMaxOpenConns(1)
 for _,q:=range []string{
 "PRAGMA journal_mode=WAL",
 "PRAGMA busy_timeout=5000",
 "CREATE TABLE IF NOT EXISTS samples(name TEXT NOT NULL,ts INTEGER NOT NULL,cpu REAL NOT NULL,memory REAL NOT NULL,disk REAL NOT NULL,payload TEXT NOT NULL,PRIMARY KEY(name,ts))",
 "CREATE INDEX IF NOT EXISTS idx_samples_ts ON samples(ts)",
 "CREATE TABLE IF NOT EXISTS traffic_daily(name TEXT NOT NULL, day TEXT NOT NULL, rx INTEGER NOT NULL DEFAULT 0, tx INTEGER NOT NULL DEFAULT 0, PRIMARY KEY(name,day))",
 "CREATE TABLE IF NOT EXISTS agent_tokens(name TEXT PRIMARY KEY,hash TEXT NOT NULL)",
 "CREATE TABLE IF NOT EXISTS node_metadata(name TEXT PRIMARY KEY,display_name TEXT NOT NULL DEFAULT '',group_name TEXT NOT NULL DEFAULT '',notes TEXT NOT NULL DEFAULT '')",
 }{if _,err=db.Exec(q);err!=nil{db.Close();return nil,fmt.Errorf("schema: %w",err)}}
 return db,nil
}
func validNodeToken(ctx context.Context,db *sql.DB,name,token,legacy string)bool{
 var h string
 err:=db.QueryRowContext(ctx,"SELECT hash FROM agent_tokens WHERE name=?",name).Scan(&h)
 if err==sql.ErrNoRows{return subtle.ConstantTimeCompare([]byte(token),[]byte(legacy))==1}
 if err!=nil{return false}
 hash:=sha256.Sum256([]byte(token));return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(hash[:])),[]byte(h))==1
}
func recordSample(db *sql.DB,s Sample)error{
 b,err:=json.Marshal(s);if err!=nil{return err}
 _,err=db.Exec("INSERT OR REPLACE INTO samples(name,ts,cpu,memory,disk,payload) VALUES(?,?,?,?,?,?)",s.Name,s.Timestamp.Unix(),s.CPU,s.Memory,s.Disk,string(b));return err
}
func queryHistory(ctx context.Context,db *sql.DB,name string,hours int)([]Point,error){
 // Fixed maximum ~300 points; group buckets for range.
 step:=int64(hours*3600/240);if step<5{step=5}
 rows,err:=db.QueryContext(ctx,"SELECT (ts / ?) * ? AS bucket,AVG(cpu),AVG(memory),AVG(disk) FROM samples WHERE name=? AND ts>=? GROUP BY bucket ORDER BY bucket ASC LIMIT 300",step,step,name,time.Now().Add(-time.Duration(hours)*time.Hour).Unix())
 if err!=nil{return nil,err};defer rows.Close()
 out:=[]Point{}
 for rows.Next(){var p Point;if err=rows.Scan(&p.Time,&p.CPU,&p.Memory,&p.Disk);err!=nil{return nil,err};out=append(out,p)}
 return out,rows.Err()
}
func restoreNodes(s *Store)error{
 rows,err:=s.db.Query("SELECT name,payload FROM samples WHERE (name,ts) IN (SELECT name,MAX(ts) FROM samples GROUP BY name) AND name NOT IN (SELECT name FROM agent_tokens WHERE hash='REVOKED')");if err!=nil{return err};defer rows.Close()
 for rows.Next(){var name,raw string;if err=rows.Scan(&name,&raw);err!=nil{return err};var p Sample;if json.Unmarshal([]byte(raw),&p)==nil{p.Name=name;s.Nodes[name]=Node{Sample:p,LastSeen:p.Timestamp}}}
 if err:=rows.Err();err!=nil{return err}
 revoked,err:=s.db.Query("SELECT name FROM agent_tokens WHERE hash='REVOKED'");if err!=nil{return err};defer revoked.Close()
 for revoked.Next(){var name string;if err:=revoked.Scan(&name);err!=nil{return err};delete(s.Nodes,name)}
 return revoked.Err()
}

type TrafficBucket struct {Time int64 `json:"time"`;RX uint64 `json:"rx"`;TX uint64 `json:"tx"`}
type TrafficSummary struct {RX uint64 `json:"rx"`;TX uint64 `json:"tx"`;Buckets []TrafficBucket `json:"buckets"`}
func validNodeName(n string)bool{if len(n)<1||len(n)>64{return false};for _,r:=range n{if !((r>='a'&&r<='z')||(r>='A'&&r<='Z')||(r>='0'&&r<='9')||r=='-'||r=='_'){return false}};return true}
func queryTraffic(ctx context.Context,db *sql.DB,name string,hours int)(TrafficSummary,error){
 // Counters are monotonic between restarts; a decrease indicates counter reset.
 // Reject implausibly large jumps; aggregate by hourly buckets (30d max 720).
 rows,err:=db.QueryContext(ctx,"SELECT ts,payload FROM samples WHERE name=? AND ts>=? ORDER BY ts",name,time.Now().Add(-time.Duration(hours)*time.Hour).Unix())
 out:=TrafficSummary{Buckets:[]TrafficBucket{}};if err!=nil{return out,err};defer rows.Close()
 buckets:=map[int64]*TrafficBucket{};var previous Sample;var seen bool
 for rows.Next(){
  var ts int64;var raw string;if err=rows.Scan(&ts,&raw);err!=nil{return out,err}
  var current Sample;if json.Unmarshal([]byte(raw),&current)!=nil{continue}
  if seen{
   var rx,tx uint64
   // Network counters can reset on reboots, or an interface can disappear.
   if current.Uptime>=previous.Uptime && current.RxBytes>=previous.RxBytes {rx=current.RxBytes-previous.RxBytes}
   if current.Uptime>=previous.Uptime && current.TxBytes>=previous.TxBytes {tx=current.TxBytes-previous.TxBytes}
   bucket:=ts/3600*3600
   b:=buckets[bucket];if b==nil{b=&TrafficBucket{Time:bucket};buckets[bucket]=b}
   b.RX+=rx;b.TX+=tx;out.RX+=rx;out.TX+=tx
  }
  previous=current;seen=true
 }
 if err=rows.Err();err!=nil{return out,err}
 keys:=make([]int64,0,len(buckets));for t:=range buckets{keys=append(keys,t)};sort.Slice(keys,func(i,j int)bool{return keys[i]<keys[j]});for _,t:=range keys{out.Buckets=append(out.Buckets,*buckets[t])}
 return out,nil
}

func sameOrigin(r *http.Request)bool{origin:=r.Header.Get("Origin");return origin==""||origin=="https://"+r.Host||origin=="http://"+r.Host}
func revokeNode(ctx context.Context,db *sql.DB,name string)error{
 // A tombstone blocks fallback authentication with the legacy shared token.
 _,err:=db.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",name,"REVOKED");return err
}
func removeNode(ctx context.Context,db *sql.DB,name string)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"DELETE FROM samples WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM node_metadata WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM traffic_daily WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",name,"REVOKED");err!=nil{return err}
 return tx.Commit()
}
func renameNode(ctx context.Context,db *sql.DB,old,new string)error{
 // The old identity is revoked. A renamed node needs a new credential and agent configuration.
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 var exists int
 if err=tx.QueryRowContext(ctx,"SELECT COUNT(*) FROM agent_tokens WHERE name=?",new).Scan(&exists);err!=nil{return err};if exists>0{return fmt.Errorf("target already registered")}
 if _,err=tx.ExecContext(ctx,"UPDATE samples SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE node_metadata SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE traffic_daily SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",old,"REVOKED");err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?)",new,"REVOKED");err!=nil{return err}
 return tx.Commit()
}

type PeriodTraffic struct{TodayRX uint64 `json:"today_rx"`;TodayTX uint64 `json:"today_tx"`;MonthRX uint64 `json:"month_rx"`;MonthTX uint64 `json:"month_tx"`}
func recordTrafficIncrement(db *sql.DB,current Sample)error{
 var raw string;var ts int64
 e:=db.QueryRow("SELECT ts,payload FROM samples WHERE name=? ORDER BY ts DESC LIMIT 1",current.Name).Scan(&ts,&raw)
 if e==sql.ErrNoRows{return nil};if e!=nil{return e}
 var old Sample;if json.Unmarshal([]byte(raw),&old)!=nil{return nil}
 dt:=current.Timestamp.Unix()-ts
 if dt<=0||dt>600||current.Uptime<old.Uptime{return nil}
 if current.RxBytes<old.RxBytes||current.TxBytes<old.TxBytes{return nil}
 rx,tx:=current.RxBytes-old.RxBytes,current.TxBytes-old.TxBytes
 // Counter values are bytes; reject implausible spikes above ~25 Gbit/s.
 if rx>uint64(dt)*3200000000||tx>uint64(dt)*3200000000{return nil}
 day:=current.Timestamp.UTC().Format("2006-01-02")
 _,e=db.Exec("INSERT INTO traffic_daily(name,day,rx,tx) VALUES(?,?,?,?) ON CONFLICT(name,day) DO UPDATE SET rx=rx+excluded.rx,tx=tx+excluded.tx",current.Name,day,rx,tx)
 return e
}
func periodTraffic(ctx context.Context,db *sql.DB)(map[string]PeriodTraffic,error){
 now:=time.Now().UTC();today:=now.Format("2006-01-02");month:=now.Format("2006-01")
 rows,e:=db.QueryContext(ctx,"SELECT name,day,rx,tx FROM traffic_daily WHERE day>=?",month+"-01");if e!=nil{return nil,e};defer rows.Close()
 out:=map[string]PeriodTraffic{}
 for rows.Next(){var name,day string;var rx,tx uint64;if e=rows.Scan(&name,&day,&rx,&tx);e!=nil{return nil,e};v:=out[name];if len(day)>=7&&day[:7]==month{v.MonthRX+=rx;v.MonthTX+=tx};if day==today{v.TodayRX+=rx;v.TodayTX+=tx};out[name]=v}
 return out,rows.Err()
}
