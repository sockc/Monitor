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
 "strings"
 "time"
)
// NodeProfile stores optional plan information shown in the dashboard.
// has_ipv4/has_ipv6 are tri-state: -1 unknown, 0 no, 1 yes.
type NodeProfile struct{
 Provider string `json:"provider"`
 CountryCode string `json:"country_code"`
 PriceValue float64 `json:"price_value"`
 PriceCurrency string `json:"price_currency"`
 BillingCycle string `json:"billing_cycle"`
 PeriodStart string `json:"period_start"`
 PortMbps int `json:"port_mbps"`
 HasIPv4 int `json:"has_ipv4"`
 HasIPv6 int `json:"has_ipv6"`
}
func defaultNodeProfile() NodeProfile {return NodeProfile{PriceCurrency:"USD",BillingCycle:"year",HasIPv4:-1,HasIPv6:-1}}
func validNodeProfile(p NodeProfile)bool{
 if len([]rune(p.Provider))>80||len(p.CountryCode)!=0&&len(p.CountryCode)!=2||p.PriceValue<0||p.PriceValue>100000000||p.PriceValue!=p.PriceValue||p.PortMbps<0||p.PortMbps>1000000||p.HasIPv4< -1||p.HasIPv4>1||p.HasIPv6< -1||p.HasIPv6>1{return false}
 for _,c:=range p.CountryCode{if c<'A'||c>'Z'{return false}}
 switch p.PriceCurrency {case "USD","CNY","EUR","GBP","HKD","JPY","SGD","TWD","AUD","CAD":default:return false}
 switch p.BillingCycle {case "month","quarter","year","one_time":default:return false}
 if p.PeriodStart!=""{d,e:=time.Parse("2006-01-02",p.PeriodStart);if e!=nil||d.Format("2006-01-02")!=p.PeriodStart{return false}}
 return true
}
func readNodeProfiles(ctx context.Context,db *sql.DB)(map[string]NodeProfile,error){
 rows,e:=db.QueryContext(ctx,"SELECT name,provider,country_code,price_value,price_currency,billing_cycle,period_start,port_mbps,has_ipv4,has_ipv6 FROM node_profile");if e!=nil{return nil,e};defer rows.Close()
 out:=map[string]NodeProfile{}
 for rows.Next(){var name string;var p NodeProfile;if e=rows.Scan(&name,&p.Provider,&p.CountryCode,&p.PriceValue,&p.PriceCurrency,&p.BillingCycle,&p.PeriodStart,&p.PortMbps,&p.HasIPv4,&p.HasIPv6);e!=nil{return nil,e};out[name]=p}
 return out,rows.Err()
}
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
 "CREATE TABLE IF NOT EXISTS traffic_quarter(name TEXT NOT NULL, ts INTEGER NOT NULL, rx INTEGER NOT NULL DEFAULT 0, tx INTEGER NOT NULL DEFAULT 0, PRIMARY KEY(name,ts))",
 "CREATE TABLE IF NOT EXISTS node_limits(name TEXT PRIMARY KEY, timezone TEXT NOT NULL DEFAULT 'UTC', quota_gb REAL NOT NULL DEFAULT 0, expires_on TEXT NOT NULL DEFAULT '')",
 "CREATE TABLE IF NOT EXISTS agent_tokens(name TEXT PRIMARY KEY,hash TEXT NOT NULL)",
 "CREATE TABLE IF NOT EXISTS node_metadata(name TEXT PRIMARY KEY,display_name TEXT NOT NULL DEFAULT '',group_name TEXT NOT NULL DEFAULT '',notes TEXT NOT NULL DEFAULT '')",
 "CREATE TABLE IF NOT EXISTS node_favorites(name TEXT PRIMARY KEY, favorite INTEGER NOT NULL DEFAULT 0 CHECK(favorite IN (0,1)))",
 "CREATE TABLE IF NOT EXISTS node_geo(name TEXT PRIMARY KEY,public_ip TEXT NOT NULL DEFAULT '',auto_location TEXT NOT NULL DEFAULT '',updated_at INTEGER NOT NULL DEFAULT 0)",
 "CREATE TABLE IF NOT EXISTS node_profile(name TEXT PRIMARY KEY,provider TEXT NOT NULL DEFAULT '',country_code TEXT NOT NULL DEFAULT '',price_value REAL NOT NULL DEFAULT 0,price_currency TEXT NOT NULL DEFAULT 'USD',billing_cycle TEXT NOT NULL DEFAULT 'year',period_start TEXT NOT NULL DEFAULT '',port_mbps INTEGER NOT NULL DEFAULT 0,has_ipv4 INTEGER NOT NULL DEFAULT -1,has_ipv6 INTEGER NOT NULL DEFAULT -1)",
 }{if _,err=db.Exec(q);err!=nil{db.Close();return nil,fmt.Errorf("schema: %w",err)}}
 // Migration for existing installations: preserve all stored metadata and
 // add a user-controlled geographic location without guessing from public IP.
 rows,err:=db.Query("PRAGMA table_info(node_metadata)")
 if err!=nil{db.Close();return nil,err}
 hasLocation:=false;hasSortOrder:=false
 for rows.Next(){var cid,notnull,pk int;var name,typ string;var defaultValue sql.NullString
  if err=rows.Scan(&cid,&name,&typ,&notnull,&defaultValue,&pk);err!=nil{rows.Close();db.Close();return nil,err}
  if name=="location"{hasLocation=true};if name=="sort_order"{hasSortOrder=true}
 }
 if err=rows.Err();err!=nil{rows.Close();db.Close();return nil,err}
 rows.Close()
 if !hasLocation{if _,err=db.Exec("ALTER TABLE node_metadata ADD COLUMN location TEXT NOT NULL DEFAULT ''");err!=nil{db.Close();return nil,fmt.Errorf("location migration: %w",err)}}
 if !hasSortOrder{if _,err=db.Exec("ALTER TABLE node_metadata ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0");err!=nil{db.Close();return nil,fmt.Errorf("sort migration: %w",err)}}
 // Older IP-location tables may not have a country code for the flag.
 geoRows,eGeo:=db.Query("PRAGMA table_info(node_geo)")
 if eGeo!=nil{db.Close();return nil,eGeo}
 hasCountry:=false
 for geoRows.Next(){var cid,notnull,pk int;var name,typ string;var def sql.NullString;if eGeo=geoRows.Scan(&cid,&name,&typ,&notnull,&def,&pk);eGeo!=nil{geoRows.Close();db.Close();return nil,eGeo};if name=="country_code"{hasCountry=true}}
 if eGeo=geoRows.Err();eGeo!=nil{geoRows.Close();db.Close();return nil,eGeo};geoRows.Close()
 if !hasCountry{if _,eGeo=db.Exec("ALTER TABLE node_geo ADD COLUMN country_code TEXT NOT NULL DEFAULT ''");eGeo!=nil{db.Close();return nil,fmt.Errorf("geo migration: %w",eGeo)}}

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
 if _,err=tx.ExecContext(ctx,"DELETE FROM node_favorites WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM node_geo WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM node_profile WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM traffic_daily WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM traffic_quarter WHERE name=?",name);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM node_limits WHERE name=?",name);err!=nil{return err}
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
 if _,err=tx.ExecContext(ctx,"UPDATE node_favorites SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE node_geo SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE node_profile SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE traffic_daily SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE traffic_quarter SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE node_limits SET name=? WHERE name=?",new,old);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?) ON CONFLICT(name) DO UPDATE SET hash=excluded.hash",old,"REVOKED");err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO agent_tokens(name,hash) VALUES(?,?)",new,"REVOKED");err!=nil{return err}
 return tx.Commit()
}

// Traffic is retained in quarter-hour UTC buckets so per-node timezone changes
// do not modify historical data. A bucket is attributed to its ending sample.
type NodeLimits struct {
 Name string `json:"name"`
 Timezone string `json:"timezone"`
 QuotaGB float64 `json:"quota_gb"`
 ExpiresOn string `json:"expires_on"`
}
func validNodeLimits(v NodeLimits)bool{
 if !validNodeName(v.Name)||v.QuotaGB<0||v.QuotaGB>1000000||v.QuotaGB!=v.QuotaGB{return false}
 if len(v.Timezone)<1||len(v.Timezone)>64||strings.Contains(v.Timezone,"..")||strings.HasPrefix(v.Timezone,"/"){return false}
 for _,c:=range v.Timezone{if !((c>='a'&&c<='z')||(c>='A'&&c<='Z')||(c>='0'&&c<='9')||c=='/'||c=='_'||c=='+'||c=='-'){return false}}
 if _,e:=time.LoadLocation(v.Timezone);e!=nil{return false}
 if v.ExpiresOn!=""{d,e:=time.Parse("2006-01-02",v.ExpiresOn);if e!=nil||d.Format("2006-01-02")!=v.ExpiresOn{return false}}
 return true
}
func setNodeLimits(ctx context.Context,db *sql.DB,v NodeLimits)error{
 _,e:=db.ExecContext(ctx,"INSERT INTO node_limits(name,timezone,quota_gb,expires_on) VALUES(?,?,?,?) ON CONFLICT(name) DO UPDATE SET timezone=excluded.timezone,quota_gb=excluded.quota_gb,expires_on=excluded.expires_on",v.Name,v.Timezone,v.QuotaGB,v.ExpiresOn)
 return e
}
func readNodeLimits(ctx context.Context,db *sql.DB)(map[string]NodeLimits,error){
 rows,e:=db.QueryContext(ctx,"SELECT name,timezone,quota_gb,expires_on FROM node_limits");if e!=nil{return nil,e};defer rows.Close()
 out:=map[string]NodeLimits{}
 for rows.Next(){var v NodeLimits;if e=rows.Scan(&v.Name,&v.Timezone,&v.QuotaGB,&v.ExpiresOn);e!=nil{return nil,e};out[v.Name]=v}
 return out,rows.Err()
}
type PeriodTraffic struct{
 TodayRX uint64 `json:"today_rx"`
 TodayTX uint64 `json:"today_tx"`
 MonthRX uint64 `json:"month_rx"`
 MonthTX uint64 `json:"month_tx"`
 Timezone string `json:"timezone"`
 QuotaGB float64 `json:"quota_gb"`
 ExpiresOn string `json:"expires_on"`
 HasSamples bool `json:"has_samples"`
 ObservedBuckets int `json:"observed_buckets"`
}
func recordTrafficIncrement(db *sql.DB,current Sample)error{
 var raw string;var ts int64
 e:=db.QueryRow("SELECT ts,payload FROM samples WHERE name=? ORDER BY ts DESC LIMIT 1",current.Name).Scan(&ts,&raw)
 if e==sql.ErrNoRows{return nil};if e!=nil{return e}
 var old Sample;if json.Unmarshal([]byte(raw),&old)!=nil{return nil}
 dt:=current.Timestamp.Unix()-ts
 if dt<=0||dt>600||current.Uptime<old.Uptime{return nil}
 if current.RxBytes<old.RxBytes||current.TxBytes<old.TxBytes{return nil}
 rx,tx:=current.RxBytes-old.RxBytes,current.TxBytes-old.TxBytes
 // Reject implausible spikes (>25.6 Gbit/s sustained over the interval).
 if rx>uint64(dt)*3200000000||tx>uint64(dt)*3200000000{return nil}
 txdb,e:=db.Begin();if e!=nil{return e};defer txdb.Rollback()
 bucket:=current.Timestamp.Unix()/900*900
 if _,e=txdb.Exec("INSERT INTO traffic_quarter(name,ts,rx,tx) VALUES(?,?,?,?) ON CONFLICT(name,ts) DO UPDATE SET rx=rx+excluded.rx,tx=tx+excluded.tx",current.Name,bucket,rx,tx);e!=nil{return e}
 day:=current.Timestamp.UTC().Format("2006-01-02")
 if _,e=txdb.Exec("INSERT INTO traffic_daily(name,day,rx,tx) VALUES(?,?,?,?) ON CONFLICT(name,day) DO UPDATE SET rx=rx+excluded.rx,tx=tx+excluded.tx",current.Name,day,rx,tx);e!=nil{return e}
 return txdb.Commit()
}
func periodTraffic(ctx context.Context,db *sql.DB)(map[string]PeriodTraffic,error){
 cfg,e:=readNodeLimits(ctx,db);if e!=nil{return nil,e}
 rows,e:=db.QueryContext(ctx,"SELECT DISTINCT name FROM traffic_quarter");if e!=nil{return nil,e}
 names:=map[string]bool{};for name:=range cfg{names[name]=true}
 for rows.Next(){var name string;if e=rows.Scan(&name);e!=nil{rows.Close();return nil,e};names[name]=true}
 if e=rows.Err();e!=nil{rows.Close();return nil,e};rows.Close()
 out:=map[string]PeriodTraffic{}
 for name:=range names{
  v:=cfg[name];if v.Timezone==""{v.Timezone="UTC"}
  loc,e:=time.LoadLocation(v.Timezone);if e!=nil{loc=time.UTC;v.Timezone="UTC"}
  local:=time.Now().In(loc)
  startToday:=time.Date(local.Year(),local.Month(),local.Day(),0,0,0,0,loc).Unix()
  startMonth:=time.Date(local.Year(),local.Month(),1,0,0,0,0,loc).Unix()
  nextMonth:=time.Date(local.Year(),local.Month()+1,1,0,0,0,0,loc).Unix()
  p:=PeriodTraffic{Timezone:v.Timezone,QuotaGB:v.QuotaGB,ExpiresOn:v.ExpiresOn}
  var rx,tx uint64
  e=db.QueryRowContext(ctx,"SELECT COALESCE(SUM(rx),0),COALESCE(SUM(tx),0),COUNT(*) FROM traffic_quarter WHERE name=? AND ts>=? AND ts<?",name,startMonth,nextMonth).Scan(&rx,&tx,&p.ObservedBuckets)
  if e!=nil{return nil,e};p.MonthRX=rx;p.MonthTX=tx;p.HasSamples=p.ObservedBuckets>0
  if e=db.QueryRowContext(ctx,"SELECT COALESCE(SUM(rx),0),COALESCE(SUM(tx),0) FROM traffic_quarter WHERE name=? AND ts>=? AND ts<?",name,startToday,nextMonth).Scan(&p.TodayRX,&p.TodayTX);e!=nil{return nil,e}
  out[name]=p
 }
 return out,nil
}
