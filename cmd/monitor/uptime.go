package main

import (
 "context"
 "database/sql"
 "math"
 "time"
)

// The uptime number is heartbeat-minute coverage at the Monitor server,
// not an ICMP or independent infrastructure SLA measurement.
type UptimeWindow struct {
 Hours int `json:"hours"`
 Percent *float64 `json:"percent"`
 ObservedMinutes int64 `json:"observed_minutes"`
 ExpectedMinutes int64 `json:"expected_minutes"`
}

func queryUptime(ctx context.Context, db *sql.DB, name string, now time.Time)([]UptimeWindow,error){
 var first sql.NullInt64
 if err:=db.QueryRowContext(ctx,"SELECT MIN(ts) FROM samples WHERE name=?",name).Scan(&first);err!=nil{return nil,err}
 out:=make([]UptimeWindow,0,3)
 currentMinute:=now.Unix()/60
 for _,hours:=range []int{24,168,720}{
  w:=UptimeWindow{Hours:hours}
  if first.Valid {
   cutoff:=now.Add(-time.Duration(hours)*time.Hour).Unix()
   from:=cutoff
   if first.Int64>from {from=first.Int64}
   // Bin by whole UTC minutes. Current minute is included.
   firstBucket:=from/60
   if currentMinute>=firstBucket{
    w.ExpectedMinutes=currentMinute-firstBucket+1
    var count int64
    err:=db.QueryRowContext(ctx,
     "SELECT COUNT(DISTINCT(ts/60)) FROM samples WHERE name=? AND ts>=? AND ts<=?",
     name,firstBucket*60,now.Unix()).Scan(&count)
    if err!=nil{return nil,err}
    w.ObservedMinutes=count
    if w.ObservedMinutes>w.ExpectedMinutes{w.ObservedMinutes=w.ExpectedMinutes}
    if w.ExpectedMinutes>=3 {
     v:=math.Round(float64(w.ObservedMinutes)/float64(w.ExpectedMinutes)*1000)/10
     w.Percent=&v
    }
   }
  }
  out=append(out,w)
 }
 return out,nil
}
