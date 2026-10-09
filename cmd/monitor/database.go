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
 "path/filepath"
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
 "CREATE TABLE IF NOT EXISTS agent_tokens(name TEXT PRIMARY KEY,hash TEXT NOT NULL)",
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
 rows,err:=s.db.Query("SELECT payload FROM samples WHERE (name,ts) IN (SELECT name,MAX(ts) FROM samples GROUP BY name)");if err!=nil{return err};defer rows.Close()
 for rows.Next(){var raw string;if err=rows.Scan(&raw);err!=nil{return err};var p Sample;if json.Unmarshal([]byte(raw),&p)==nil{s.Nodes[p.Name]=Node{Sample:p,LastSeen:p.Timestamp}}}
 return rows.Err()
}
