package main

import (
 "strings"
 "testing"
)

func TestNetifdWAN(t *testing.T) {
 tests:=[]struct{input,device string;ok bool}{
  {`{"up":true,"device":"eth0.2","l3_device":"pppoe-wan"}`,"pppoe-wan",true},
  {`{"up":true,"device":"eth0.2"}`,"eth0.2",true},
  {`{"up":false,"device":"eth0"}`,"",false},
  {`{"up":true,"l3_device":"../../etc/passwd"}`,"",false},
 }
 for _,tc:=range tests {device,ok:=parseIfStatus([]byte(tc.input));if device!=tc.device||ok!=tc.ok {t.Errorf("parse %q got %q %v, want %q %v",tc.input,device,ok,tc.device,tc.ok)}}
}
func TestDefaultRouteSkipsTUN(t *testing.T) {
 raw:="Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT\n"+
  "tun0 00000000 0100000A 0003 0 0 0 00000000 0 0 0\n"+
  "pppoe-wan 00000000 0100000A 0003 0 0 100 00000000 0 0 0\n"
 if got:=fallbackDefaultWAN(raw);got!="pppoe-wan"{t.Fatalf("got %q",got)}
 if got:=fallbackDefaultWAN(strings.Replace(raw,"pppoe-wan","wg0",1));got!=""{t.Fatalf("tunnel used %q",got)}
}
func TestTemperatureParsing(t *testing.T){
 for _,tt:=range []struct{raw string;want float64;valid bool}{
  {"53456\n",53.456,true},{"42",42,true},{"-10",-10,false},{"400000",400,false},{"bad",0,false},
 } {v,ok:=parseTemp(tt.raw);if ok!=tt.valid||(ok&&v!=tt.want){t.Errorf("input=%q got=%v %v",tt.raw,v,ok)}}
}
