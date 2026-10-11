package main

import (
 "context"
 "encoding/json"
 "os"
 "os/exec"
 "strconv"
 "strings"
 "sync"
 "time"
)

// OpenWrt netifd is authoritative for WAN identity. A default route may
// lead through a VPN/TUN device and must not be used as the first choice.
type ifStatus struct {
 Up bool `json:"up"`
 Device string `json:"device"`
 L3Device string `json:"l3_device"`
}
func parseIfStatus(raw []byte) (string, bool) {
 var s ifStatus
 if json.Unmarshal(raw, &s) != nil || !s.Up { return "", false }
 for _, name := range []string{s.L3Device, s.Device} {
  if validNetworkDevice(name) { return name, true }
 }
 return "", false
}
func validNetworkDevice(name string) bool {
 if name == "" || len(name) > 32 { return false }
 for _, c := range name {
  if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' && c != '-' && c != '.' { return false }
 }
 return name != "lo"
}
func netifdDevice(logical string) string {
 ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
 defer cancel()
 output, err := exec.CommandContext(ctx, "ubus", "call", "network.interface."+logical, "status").Output()
 if err != nil || len(output) > 32*1024 { return "" }
 name, ok := parseIfStatus(output)
 if !ok { return "" }
 if _, err := os.Stat("/sys/class/net/"+name); err != nil { return "" }
 return name
}
func fallbackDefaultWAN(data string) string {
 for index, line := range strings.Split(data, "\n") {
  if index == 0 { continue }
  f := strings.Fields(line)
  if len(f) < 4 || f[1] != "00000000" || !validNetworkDevice(f[0]) { continue }
  flags, err := strconv.ParseUint(f[3], 16, 32)
  if err != nil || flags&1 == 0 { continue }
  lower := strings.ToLower(f[0])
  if strings.HasPrefix(lower, "tun") || strings.HasPrefix(lower, "wg") ||
    strings.HasPrefix(lower, "tailscale") || strings.HasPrefix(lower, "clash") ||
    strings.HasPrefix(lower, "utun") || strings.HasPrefix(lower, "tap") { continue }
  return f[0]
 }
 return ""
}
func findWAN() (string,string) {
 wan := netifdDevice("wan")
 wan6 := netifdDevice("wan6")
 if wan == "" {
  // IPv6-only setups are valid. Never sum WAN and WAN6 counters:
  // both commonly resolve to the same lower network device.
  if wan6 != "" { wan = wan6
  } else if raw, err := os.ReadFile("/proc/net/route"); err == nil {
   wan = fallbackDefaultWAN(string(raw))
  }
 }
 return wan, wan6
}
type wanCache struct {
 mu sync.Mutex
 device, wan6 string
 next time.Time
}
func (c *wanCache) get() (string,string) {
 c.mu.Lock()
 defer c.mu.Unlock()
 if time.Now().Before(c.next) { return c.device,c.wan6 }
 c.device,c.wan6=findWAN()
 c.next=time.Now().Add(60*time.Second)
 return c.device,c.wan6
}

// Temperature is optional. Reporting an unsupported sensor as 0C would
// incorrectly imply a healthy, meaningful measurement.
func parseTemp(raw string) (float64,bool) {
 v,err:=strconv.ParseFloat(strings.TrimSpace(raw),64)
 if err!=nil {return 0,false}
 if v>1000 {v/=1000}
 if v<=0||v>125 {return 0,false}
 return v,true
}
func cpuThermal() *float64 {
 zones,_:=os.ReadDir("/sys/class/thermal")
 for _,z:=range zones {
  if !strings.HasPrefix(z.Name(),"thermal_zone") {continue}
  path:="/sys/class/thermal/"+z.Name()
  kind,_:=os.ReadFile(path+"/type")
  name:=strings.ToLower(strings.TrimSpace(string(kind)))
  if !strings.Contains(name,"cpu")&&!strings.Contains(name,"soc")&&!strings.Contains(name,"package")&&!strings.Contains(name,"processor") {continue}
  b,e:=os.ReadFile(path+"/temp");if e!=nil {continue}
  if v,ok:=parseTemp(string(b));ok{return &v}
 }
 hwmons,_:=os.ReadDir("/sys/class/hwmon")
 for _,h:=range hwmons{
  if !strings.HasPrefix(h.Name(),"hwmon"){continue}
  path:="/sys/class/hwmon/"+h.Name()
  kind,_:=os.ReadFile(path+"/name")
  label,_:=os.ReadFile(path+"/temp1_label")
  desc:=strings.ToLower(string(kind)+" "+string(label))
  if !strings.Contains(desc,"cpu")&&!strings.Contains(desc,"soc")&&!strings.Contains(desc,"coretemp"){continue}
  b,e:=os.ReadFile(path+"/temp1_input");if e!=nil {continue}
  if v,ok:=parseTemp(string(b));ok{return &v}
 }
 return nil
}
