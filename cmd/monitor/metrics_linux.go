//go:build linux

package main

import (
 "bufio"
 "os"
 "runtime"
 "strconv"
 "strings"
 "syscall"
 "time"
)
type cpuTimes struct{idle,total uint64}
func readCPU()cpuTimes{
 b,e:=os.ReadFile("/proc/stat");if e!=nil{return cpuTimes{}}
 line:=strings.SplitN(string(b),"\n",2)[0];f:=strings.Fields(line);if len(f)<5{return cpuTimes{}}
 var t cpuTimes
 for i,x:=range f[1:]{v,_:=strconv.ParseUint(x,10,64);t.total+=v;if i==3||i==4{t.idle+=v}}
 return t
}
func percent(used,total uint64)float64{if total==0{return 0};return float64(used)*100/float64(total)}
func collect(name string,prev cpuTimes)(Sample,cpuTimes){
 now:=readCPU();var cpu float64
 if now.total>prev.total&&now.idle>=prev.idle{cpu=percent(now.total-prev.total-(now.idle-prev.idle),now.total-prev.total)}
 host,_:=os.Hostname()
 s:=Sample{Name:name,Hostname:host,OS:linuxDistribution(),Arch:runtime.GOARCH,CPUCores:runtime.NumCPU(),CPU:cpu,Timestamp:time.Now().UTC()}
 f,e:=os.Open("/proc/meminfo");if e==nil{
  m:=map[string]uint64{};sc:=bufio.NewScanner(f)
  for sc.Scan(){parts:=strings.Fields(sc.Text());if len(parts)>=2{v,_:=strconv.ParseUint(parts[1],10,64);m[strings.TrimSuffix(parts[0],":")]=v}}
  f.Close();if m["MemTotal"]>0{total:=m["MemTotal"]*1024;available:=m["MemAvailable"]*1024;if available>total{available=total};s.MemoryTotal=total;s.MemoryUsed=total-available;s.Memory=percent(s.MemoryUsed,total)}
 }
 var st syscall.Statfs_t
 if syscall.Statfs("/",&st)==nil&&st.Blocks>0{total:=st.Blocks*uint64(st.Bsize);free:=st.Bavail*uint64(st.Bsize);if free>total{free=total};s.DiskTotal=total;s.DiskUsed=total-free;s.Disk=percent(s.DiskUsed,total)}
 b,e:=os.ReadFile("/proc/net/dev");if e==nil{
  for _,line:=range strings.Split(string(b),"\n"){parts:=strings.SplitN(line,":",2);if len(parts)!=2||strings.TrimSpace(parts[0])=="lo"{continue};v:=strings.Fields(parts[1]);if len(v)<16{continue};rx,_:=strconv.ParseUint(v[0],10,64);tx,_:=strconv.ParseUint(v[8],10,64);s.RxBytes+=rx;s.TxBytes+=tx}
 }
 b,e=os.ReadFile("/proc/uptime");if e==nil{v:=strings.Fields(string(b));if len(v)>0{n,_:=strconv.ParseFloat(v[0],64);s.Uptime=uint64(n)}}
 return s,now
}

func linuxDistribution()string{
 b,e:=os.ReadFile("/etc/os-release");if e!=nil{return "Linux"}
 name,version,pretty:="","",""
 for _,line:=range strings.Split(string(b),"\n"){
  if !strings.Contains(line,"="){continue}
  pair:=strings.SplitN(line,"=",2);value:=strings.Trim(pair[1], " \"'")
  switch pair[0]{
  case "NAME":name=value
  case "VERSION_ID":version=value
  case "PRETTY_NAME":pretty=value
  }
 }
 if name!=""{if version!=""{return name+" "+version};return name}
 if pretty!=""{return pretty}
 return "Linux"
}
