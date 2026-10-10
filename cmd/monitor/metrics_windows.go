//go:build windows

package main

import (
 "os"
 "os/exec"
 "runtime"
 "strconv"
 "strings"
 "sync"
 "syscall"
 "time"
 "unsafe"
)

// Native Win32 counters are sampled without PowerShell or WMI subprocesses.
// netstat -e is used for basic total Ethernet counters; per-adapter accounting
// is a separate feature and should not be inferred from the aggregate here.
var (
 kernel32=syscall.NewLazyDLL("kernel32.dll")
 getSystemTimes=kernel32.NewProc("GetSystemTimes")
 globalMemoryStatusEx=kernel32.NewProc("GlobalMemoryStatusEx")
 getDiskFreeSpaceExW=kernel32.NewProc("GetDiskFreeSpaceExW")
 getTickCount64=kernel32.NewProc("GetTickCount64")
 cpuModelOnce sync.Once
 cachedCPUModel string
)
type winFiletime struct{Low,High uint32}
func (f winFiletime) value()uint64{return uint64(f.High)<<32|uint64(f.Low)}
type memoryStatusEx struct {
 Length,Load uint32
 TotalPhys,AvailPhys,TotalPageFile,AvailPageFile,TotalVirtual,AvailVirtual,AvailExtendedVirtual uint64
}
type cpuTimes struct{idle,total uint64}
func readCPU()cpuTimes {
 var idle,kernel,user winFiletime
 ok,_,_:=getSystemTimes.Call(uintptr(unsafe.Pointer(&idle)),uintptr(unsafe.Pointer(&kernel)),uintptr(unsafe.Pointer(&user)))
 if ok==0{return cpuTimes{}}
 return cpuTimes{idle:idle.value(),total:kernel.value()+user.value()}
}
func winPercent(used,total uint64)float64{
 if total==0{return 0}
 if used>total{used=total}
 return float64(used)*100/float64(total)
}
func windowsCPUModel()string {
 cpuModelOnce.Do(func(){
  out,e:=exec.Command("reg.exe","query",`HKLM\HARDWARE\DESCRIPTION\System\CentralProcessor\0`,"/v","ProcessorNameString").Output()
  if e!=nil{return}
  for _,line:=range strings.Split(string(out),"\n"){
   if i:=strings.Index(line,"REG_SZ");i>=0{
    cachedCPUModel=strings.TrimSpace(line[i+len("REG_SZ"):])
    if len(cachedCPUModel)>120{cachedCPUModel=cachedCPUModel[:120]}
    return
   }
  }
 })
 return cachedCPUModel
}
func windowsDisk()(total,used uint64){
 drive:=os.Getenv("SystemDrive");if len(drive)!=2||drive[1]!=':'{drive="C:"}
 root:=drive+"\\"
 p,e:=syscall.UTF16PtrFromString(root);if e!=nil{return 0,0}
 var available,capacity,free uint64
 ok,_,_:=getDiskFreeSpaceExW.Call(uintptr(unsafe.Pointer(p)),uintptr(unsafe.Pointer(&available)),uintptr(unsafe.Pointer(&capacity)),uintptr(unsafe.Pointer(&free)))
 if ok==0||capacity==0{return 0,0}
 if free>capacity{free=capacity}
 return capacity,capacity-free
}
func windowsNetwork()(received,sent uint64){
 // Locale-independent detection: the first data row with two decimal values
 // is the Ethernet byte counter in Windows' "netstat -e" output.
 out,e:=exec.Command("netstat.exe","-e").Output()
 if e!=nil{return 0,0}
 for _,line:=range strings.Split(string(out),"\n"){
  fields:=strings.Fields(line)
  if len(fields)<3{continue}
  a,errA:=strconv.ParseUint(fields[len(fields)-2],10,64)
  b,errB:=strconv.ParseUint(fields[len(fields)-1],10,64)
  if errA==nil&&errB==nil{return a,b}
 }
 return 0,0
}
func collect(name string,prev cpuTimes)(Sample,cpuTimes){
 now:=readCPU()
 host,_:=os.Hostname()
 sample:=Sample{Name:name,Hostname:host,OS:"Windows",Arch:runtime.GOARCH,CPUCores:runtime.NumCPU(),CPUModel:windowsCPUModel(),Timestamp:time.Now().UTC()}
 if now.total>prev.total&&now.idle>=prev.idle&&now.idle-prev.idle<=now.total-prev.total{
  sample.CPU=winPercent(now.total-prev.total-(now.idle-prev.idle),now.total-prev.total)
 }
 var mem memoryStatusEx
 mem.Length=uint32(unsafe.Sizeof(mem))
 if ok,_,_:=globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)));ok!=0&&mem.TotalPhys>0{
  sample.MemoryTotal=mem.TotalPhys
  available:=mem.AvailPhys;if available>mem.TotalPhys{available=mem.TotalPhys}
  sample.MemoryUsed=mem.TotalPhys-available
  sample.Memory=winPercent(sample.MemoryUsed,sample.MemoryTotal)
 }
 sample.DiskTotal,sample.DiskUsed=windowsDisk()
 sample.Disk=winPercent(sample.DiskUsed,sample.DiskTotal)
 sample.RxBytes,sample.TxBytes=windowsNetwork()
 if up,_,_:=getTickCount64.Call();up>0{sample.Uptime=uint64(up)/1000}
 return sample,now
}
