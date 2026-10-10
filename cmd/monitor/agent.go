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
 "sync"
 "time"
)
func runAgent(server,name,token string,interval time.Duration){runAgentContext(context.Background(),server,name,token,interval)}
func runAgentContext(parent context.Context,server,name,token string,interval time.Duration){
 if server==""{server=os.Getenv("MONITOR_SERVER")}
 if !strings.HasPrefix(server,"https://")&&!strings.HasPrefix(server,"http://127.0.0.1:"){log.Fatal("MONITOR_SERVER must use HTTPS (or localhost for testing)")}
 if name==""{name=os.Getenv("MONITOR_NODE_NAME")};if name==""{name,_=os.Hostname()};if name==""{log.Fatal("name is required")}
 if interval<time.Second{log.Fatal("interval must be >= 1s")}
 client:=&http.Client{Timeout:10*time.Second}
 prev:=readCPU()
 geoEnabled:=os.Getenv("MONITOR_GEO_ENABLED")!="0"
 ipEnabled:=os.Getenv("MONITOR_IP_DETECTION_ENABLED")!="0"
 type netSnapshot struct{geo agentGeo;ipv4 string;ipv6 string}
 var mu sync.RWMutex
 current:=netSnapshot{}
 // Outbound HTTP checks are intentionally independent of telemetry samples.
 // Avoid delaying the five-second metrics loop on slow or broken IPv6 routes.
 if geoEnabled||ipEnabled{go func(){
  for {
   mu.RLock();next:=current;mu.RUnlock()
   if geoEnabled {
    if discovered,err:=lookupAgentGeo(client);err==nil{next.geo=discovered}else{log.Printf("auto location: %v",err)}
   }
   if ipEnabled {
    next.ipv4="";next.ipv6=""
    if ip,err:=lookupOutboundIP("4");err==nil{next.ipv4=ip}else{log.Printf("IPv4 outbound detection: %v",err)}
    if ip,err:=lookupOutboundIP("6");err==nil{next.ipv6=ip}else{log.Printf("IPv6 outbound detection: %v",err)}
   }
   mu.Lock();current=next;mu.Unlock()
   time.Sleep(12*time.Hour)
  }
 }()}
 send:=func(){
  sample,next:=collect(name,prev);prev=next;sample.AgentVersion=monitorVersion
  mu.RLock();snapshot:=current;mu.RUnlock()
  sample.PublicIP=snapshot.geo.IP;sample.AutoLocation=snapshot.geo.Location;sample.CountryCode=snapshot.geo.CountryCode
  sample.PublicIPv4=snapshot.ipv4;sample.PublicIPv6=snapshot.ipv6
  b,e:=json.Marshal(sample);if e!=nil{return}
  ctx,cancel:=context.WithTimeout(parent,10*time.Second);defer cancel()
  req,e:=http.NewRequestWithContext(ctx,"POST",strings.TrimRight(server,"/")+"/api/v1/ingest",bytes.NewReader(b));if e!=nil{log.Print(e);return}
  req.Header.Set("Authorization","Bearer "+token);req.Header.Set("Content-Type","application/json")
  resp,e:=client.Do(req);if e!=nil{log.Printf("report: %v",e);return};defer resp.Body.Close();io.Copy(io.Discard,io.LimitReader(resp.Body,1024))
  if resp.StatusCode!=204{log.Printf("server rejected sample: HTTP %d",resp.StatusCode)}
 }
 log.Printf("agent %s started",name);send()
 ticker:=time.NewTicker(interval);defer ticker.Stop();for {select {case <-parent.Done():return;case <-ticker.C:send()}}
}

type agentGeo struct{IP string;Location string;CountryCode string}
// The agent reports its own egress address, never the reverse proxy IP.
// Keep provider use infrequent and optional, and do not block monitoring for
// longer than the HTTP timeout when the lookup fails.
func lookupAgentGeo(c *http.Client)(agentGeo,error){
 ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second);defer cancel()
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,"https://ipwho.is/?lang=zh-CN",nil);if err!=nil{return agentGeo{},err}
 req.Header.Set("User-Agent","Monitor-Agent/"+monitorVersion)
 resp,err:=c.Do(req);if err!=nil{return agentGeo{},err};defer resp.Body.Close()
 if resp.StatusCode!=200{return agentGeo{},fmt.Errorf("location API returned HTTP %d",resp.StatusCode)}
 var result struct{Success bool `json:"success"`;IP string `json:"ip"`;Country string `json:"country"`;CountryCode string `json:"country_code"`;City string `json:"city"`;Region string `json:"region"`}
 if err=json.NewDecoder(io.LimitReader(resp.Body,16*1024)).Decode(&result);err!=nil{return agentGeo{},err}
 if !result.Success||!validPublicIP(result.IP){return agentGeo{},fmt.Errorf("location API did not return a valid public address")}
 area:=strings.TrimSpace(result.City)
 if area==""{area=strings.TrimSpace(result.Region)}
 country:=strings.TrimSpace(result.Country)
 if len([]rune(area))>60||len([]rune(country))>60{return agentGeo{},fmt.Errorf("location response too long")}
 location:=country
 if area!=""&&area!=country{if location!=""{location+=" · "};location+=area}
 if location==""{return agentGeo{},fmt.Errorf("location API returned no region")}
 countryCode:=strings.ToUpper(strings.TrimSpace(result.CountryCode));if len(countryCode)!=2{countryCode=""};for _,c:=range countryCode{if c<'A'||c>'Z'{countryCode=""}};return agentGeo{IP:result.IP,Location:location,CountryCode:countryCode},nil
}
func validPublicIP(text string)bool{
 ip:=net.ParseIP(strings.TrimSpace(text))
 return ip!=nil&&!ip.IsPrivate()&&!ip.IsLoopback()&&!ip.IsLinkLocalUnicast()&&!ip.IsLinkLocalMulticast()&&!ip.IsUnspecified()&&!ip.IsMulticast()
}


func outboundIPClient(family string)*http.Client{
 transport:=&http.Transport{
  Proxy:nil, // Do not classify an HTTP proxy's IP as the VPS interface.
  DialContext:func(ctx context.Context,network,address string)(net.Conn,error){
   dialer:=&net.Dialer{Timeout:3*time.Second}
   return dialer.DialContext(ctx,"tcp"+family,address)
  },
  TLSHandshakeTimeout:3*time.Second,
  DisableKeepAlives:true,
 }
 return &http.Client{Timeout:4*time.Second,Transport:transport}
}
// A successful IPv4/IPv6-specific HTTPS request confirms outbound access.
// Failure is treated as unknown, not proof the VPS lacks that IP family.
func lookupOutboundIP(family string)(string,error){
 client:=outboundIPClient(family)
 defer client.CloseIdleConnections()
 return lookupIPFamily(client,family)
}
func lookupIPFamily(c *http.Client,family string)(string,error){
 if family!="4"&&family!="6"{return "",fmt.Errorf("unknown IP family")}
 ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second);defer cancel()
 url:="https://api"+family+".ipify.org"
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil);if err!=nil{return "",err}
 req.Header.Set("User-Agent","Monitor-Agent/"+monitorVersion)
 resp,err:=c.Do(req);if err!=nil{return "",err};defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK{return "",fmt.Errorf("IP family %s service returned HTTP %d",family,resp.StatusCode)}
 b,err:=io.ReadAll(io.LimitReader(resp.Body,128));if err!=nil{return "",err}
 ip:=strings.TrimSpace(string(b))
 parsed:=net.ParseIP(ip)
 if parsed==nil||!validPublicIP(ip){return "",fmt.Errorf("invalid public IPv%s response",family)}
 if (family=="4")!=(parsed.To4()!=nil){return "",fmt.Errorf("unexpected IP family")}
 return ip,nil
}
