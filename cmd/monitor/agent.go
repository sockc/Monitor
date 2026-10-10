package main

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net"
 "log"
 "net/http"
 "os"
 "strings"
 "time"
)
func runAgent(server,name,token string,interval time.Duration){
 if server==""{server=os.Getenv("MONITOR_SERVER")}
 if !strings.HasPrefix(server,"https://")&&!strings.HasPrefix(server,"http://127.0.0.1:"){log.Fatal("MONITOR_SERVER must use HTTPS (or localhost for testing)")}
 if name==""{name=os.Getenv("MONITOR_NODE_NAME")};if name==""{name,_=os.Hostname()};if name==""{log.Fatal("name is required")}
 if interval<time.Second{log.Fatal("interval must be >= 1s")}
 client:=&http.Client{Timeout:10*time.Second}
 prev:=readCPU()
 geoEnabled:=os.Getenv("MONITOR_GEO_ENABLED")!="0"
 var geo agentGeo
 var nextGeo time.Time
 send:=func(){
  if geoEnabled&&time.Now().After(nextGeo){
   if discovered,err:=lookupAgentGeo(client);err==nil{geo=discovered;nextGeo=time.Now().Add(12*time.Hour)}else{nextGeo=time.Now().Add(time.Hour);log.Printf("auto location: %v",err)}
  }
  sample,next:=collect(name,prev);prev=next;sample.AgentVersion=monitorVersion;sample.PublicIP=geo.IP;sample.AutoLocation=geo.Location
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

type agentGeo struct{IP string;Location string}
// The agent reports its own egress address, never the reverse proxy IP.
// Keep provider use infrequent and optional, and do not block monitoring for
// longer than the HTTP timeout when the lookup fails.
func lookupAgentGeo(c *http.Client)(agentGeo,error){
 ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second);defer cancel()
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,"https://ipwho.is/?lang=zh-CN",nil);if err!=nil{return agentGeo{},err}
 req.Header.Set("User-Agent","Monitor-Agent/"+monitorVersion)
 resp,err:=c.Do(req);if err!=nil{return agentGeo{},err};defer resp.Body.Close()
 if resp.StatusCode!=200{return agentGeo{},fmt.Errorf("location API returned HTTP %d",resp.StatusCode)}
 var result struct{Success bool `json:"success"`;IP string `json:"ip"`;Country string `json:"country"`;City string `json:"city"`;Region string `json:"region"`}
 if err=json.NewDecoder(io.LimitReader(resp.Body,16*1024)).Decode(&result);err!=nil{return agentGeo{},err}
 if !result.Success||!validPublicIP(result.IP){return agentGeo{},fmt.Errorf("location API did not return a valid public address")}
 area:=strings.TrimSpace(result.City)
 if area==""{area=strings.TrimSpace(result.Region)}
 country:=strings.TrimSpace(result.Country)
 if len([]rune(area))>60||len([]rune(country))>60{return agentGeo{},fmt.Errorf("location response too long")}
 location:=country
 if area!=""&&area!=country{if location!=""{location+=" · "};location+=area}
 if location==""{return agentGeo{},fmt.Errorf("location API returned no region")}
 return agentGeo{IP:result.IP,Location:location},nil
}
func validPublicIP(text string)bool{
 ip:=net.ParseIP(strings.TrimSpace(text))
 return ip!=nil&&!ip.IsPrivate()&&!ip.IsLoopback()&&!ip.IsLinkLocalUnicast()&&!ip.IsLinkLocalMulticast()&&!ip.IsUnspecified()&&!ip.IsMulticast()
}

