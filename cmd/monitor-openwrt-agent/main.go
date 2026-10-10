package main

import (
 "bufio"
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "io"
 "log"
 "net/http"
 "os"
 "runtime"
 "strconv"
 "strings"
 "syscall"
 "time"
)

var version = "dev"

type config struct {
 Server string `json:"server"`
 Name string `json:"name"`
 Token string `json:"token"`
 Interval int `json:"interval_seconds"`
}
type sample struct {
 Name string `json:"name"`
 Hostname string `json:"hostname"`
 OS string `json:"os"`
 Arch string `json:"arch"`
 AgentVersion string `json:"agent_version"`
 CPUCores int `json:"cpu_cores"`
 CPU float64 `json:"cpu"`
 Memory float64 `json:"memory"`
 Disk float64 `json:"disk"`
 MemoryTotal uint64 `json:"memory_total"`
 MemoryUsed uint64 `json:"memory_used"`
 DiskTotal uint64 `json:"disk_total"`
 DiskUsed uint64 `json:"disk_used"`
 RxBytes uint64 `json:"rx_bytes"`
 TxBytes uint64 `json:"tx_bytes"`
 Uptime uint64 `json:"uptime"`
 Load1 float64 `json:"load1,omitempty"`
 Load5 float64 `json:"load5,omitempty"`
 Load15 float64 `json:"load15,omitempty"`
 Timestamp time.Time `json:"timestamp"`
}
type cpuStat struct{ idle,total uint64 }
func readCPU()cpuStat {
 b,e:=os.ReadFile("/proc/stat");if e!=nil{return cpuStat{}}
 lines:=strings.SplitN(string(b),"\n",2);f:=strings.Fields(lines[0]);if len(f)<5{return cpuStat{}}
 var s cpuStat
 for i,v:=range f[1:] {n,_:=strconv.ParseUint(v,10,64);s.total+=n;if i==3||i==4{s.idle+=n}}
 return s
}
func percent(a,b uint64)float64 {if b==0{return 0};return float64(a)*100/float64(b)}
func memInfo()map[string]uint64 {
 m:=map[string]uint64{};f,e:=os.Open("/proc/meminfo");if e!=nil{return m};defer f.Close()
 sc:=bufio.NewScanner(f);for sc.Scan(){x:=strings.Fields(sc.Text());if len(x)<2{continue};n,_:=strconv.ParseUint(x[1],10,64);m[strings.TrimSuffix(x[0],":")]=n*1024};return m
}
func wanInterface()string {
 // Select IPv4 default route; use the interface counters to avoid
 // double counting bridge, veth, VPN and physical interfaces.
 f,e:=os.Open("/proc/net/route");if e!=nil{return ""};defer f.Close()
 sc:=bufio.NewScanner(f);first:=true
 for sc.Scan(){if first{first=false;continue};x:=strings.Fields(sc.Text());if len(x)<4||x[1]!="00000000"{continue};flags,e:=strconv.ParseUint(x[3],16,32);if e==nil&&flags&1!=0{return x[0]}}
 return ""
}
func networkBytes(iface string)(uint64,uint64){
 if iface==""{return 0,0}
 b,e:=os.ReadFile("/proc/net/dev");if e!=nil{return 0,0}
 for _,line:=range strings.Split(string(b),"\n"){a:=strings.SplitN(line,":",2);if len(a)!=2||strings.TrimSpace(a[0])!=iface{continue};v:=strings.Fields(a[1]);if len(v)<16{return 0,0};rx,_:=strconv.ParseUint(v[0],10,64);tx,_:=strconv.ParseUint(v[8],10,64);return rx,tx}
 return 0,0
}
func operatingSystem()string {
 b,e:=os.ReadFile("/etc/openwrt_release");if e!=nil{return "OpenWrt"}
 for _,line:=range strings.Split(string(b),"\n"){if strings.HasPrefix(line,"DISTRIB_RELEASE="){return "OpenWrt "+strings.Trim(strings.TrimPrefix(line,"DISTRIB_RELEASE="),string([]byte{32,39,34}))}}
 return "OpenWrt"
}
func collect(name string,prev cpuStat)(sample,cpuStat){
 now:=readCPU();var cpu float64;if now.total>prev.total&&now.idle>=prev.idle{cpu=percent(now.total-prev.total-(now.idle-prev.idle),now.total-prev.total)}
 host,_:=os.Hostname();s:=sample{Name:name,Hostname:host,OS:operatingSystem(),Arch:runtime.GOARCH,AgentVersion:version,CPUCores:runtime.NumCPU(),CPU:cpu,Timestamp:time.Now().UTC()}
 mem:=memInfo();s.MemoryTotal=mem["MemTotal"];available:=mem["MemAvailable"];if available==0{available=mem["MemFree"]+mem["Buffers"]+mem["Cached"]};if available>s.MemoryTotal{available=s.MemoryTotal};s.MemoryUsed=s.MemoryTotal-available;s.Memory=percent(s.MemoryUsed,s.MemoryTotal)
 var st syscall.Statfs_t
 if syscall.Statfs("/overlay",&st)!=nil{_ = syscall.Statfs("/",&st)}
 if st.Blocks>0{size:=uint64(st.Bsize);s.DiskTotal=st.Blocks*size;free:=st.Bavail*size;if free>s.DiskTotal{free=s.DiskTotal};s.DiskUsed=s.DiskTotal-free;s.Disk=percent(s.DiskUsed,s.DiskTotal)}
 s.RxBytes,s.TxBytes=networkBytes(wanInterface())
 if b,e:=os.ReadFile("/proc/uptime");e==nil{p:=strings.Fields(string(b));if len(p)>0{v,_:=strconv.ParseFloat(p[0],64);s.Uptime=uint64(v)}}
 if b,e:=os.ReadFile("/proc/loadavg");e==nil{p:=strings.Fields(string(b));if len(p)>2{s.Load1,_=strconv.ParseFloat(p[0],64);s.Load5,_=strconv.ParseFloat(p[1],64);s.Load15,_=strconv.ParseFloat(p[2],64)}}
 return s,now
}
func readConfig(path string)(config,error) {
 var c config
 st,e:=os.Stat(path);if e!=nil{return c,e};if st.Mode().Perm()&0077!=0{return c,errors.New("config must be owner-only (chmod 600)")}
 f,e:=os.Open(path);if e!=nil{return c,e};defer f.Close()
 e=json.NewDecoder(io.LimitReader(f,8192)).Decode(&c);if e!=nil{return c,e}
 if !strings.HasPrefix(c.Server,"https://"){return c,errors.New("server URL must use https")}
 if c.Token==""||strings.TrimSpace(c.Name)==""{return c,errors.New("name and token required")}
 if c.Interval==0{c.Interval=10};if c.Interval<5||c.Interval>300{return c,errors.New("interval_seconds must be 5..300")}
 return c,nil
}
func report(client *http.Client,c config,s sample)error {
 b,e:=json.Marshal(s);if e!=nil{return e}
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 req,e:=http.NewRequestWithContext(ctx,http.MethodPost,strings.TrimRight(c.Server,"/")+"/api/v1/ingest",bytes.NewReader(b));if e!=nil{return e}
 req.Header.Set("Authorization","Bearer "+c.Token);req.Header.Set("Content-Type","application/json")
 resp,e:=client.Do(req);if e!=nil{return e};defer resp.Body.Close();io.Copy(io.Discard,io.LimitReader(resp.Body,512))
 if resp.StatusCode!=http.StatusNoContent{return fmt.Errorf("ingest HTTP %d",resp.StatusCode)};return nil
}
func main(){
 file:=flag.String("config","/etc/monitor-agent.json","owner-only JSON configuration");flag.Parse()
 c,e:=readConfig(*file);if e!=nil{log.Fatal(e)}
 client:=&http.Client{Timeout:10*time.Second}
 prev:=readCPU();ticker:=time.NewTicker(time.Duration(c.Interval)*time.Second);defer ticker.Stop()
 log.Printf("OpenWrt monitor agent %s started",version)
 for {s,next:=collect(c.Name,prev);prev=next;if e:=report(client,c,s);e!=nil{log.Printf("report: %v",e)};<-ticker.C}
}
