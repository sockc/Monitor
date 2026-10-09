package main

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "io"
 "log"
 "net/http"
 "os"
 "strings"
 "time"
)
func runAgent(server,name,token string,interval time.Duration){
 if server==""{server=os.Getenv("MONITOR_SERVER")}
 if !strings.HasPrefix(server,"https://")&&!strings.HasPrefix(server,"http://127.0.0.1:"){log.Fatal("MONITOR_SERVER must use HTTPS (or localhost for testing)")}
 if name==""{name,_=os.Hostname()};if name==""{log.Fatal("name is required")}
 if interval<time.Second{log.Fatal("interval must be >= 1s")}
 client:=&http.Client{Timeout:10*time.Second}
 prev:=readCPU()
 send:=func(){
  sample,next:=collect(name,prev);prev=next
  b,e:=json.Marshal(sample);if e!=nil{return}
  ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
  req,e:=http.NewRequestWithContext(ctx,"POST",strings.TrimRight(server,"/")+"/api/v1/ingest",bytes.NewReader(b));if e!=nil{log.Print(e);return}
  req.Header.Set("Authorization","Bearer "+token);req.Header.Set("Content-Type","application/json")
  resp,e:=client.Do(req);if e!=nil{log.Printf("report: %v",e);return};defer resp.Body.Close();io.Copy(io.Discard,io.LimitReader(resp.Body,1024))
  if resp.StatusCode!=204{log.Printf("server rejected sample: HTTP %d",resp.StatusCode)}
 }
 log.Printf("agent %s started",name);send()
 ticker:=time.NewTicker(interval);defer ticker.Stop();for range ticker.C{send()}
}
var _=fmt.Sprintf
