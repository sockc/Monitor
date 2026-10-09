package main

import (
 "context"
 "path/filepath"
 "testing"
 "time"
)

func TestTrafficPeriodsAndReset(t *testing.T) {
 db,err:=openDB(filepath.Join(t.TempDir(),"monitor.db"));if err!=nil{t.Fatal(err)};defer db.Close()
 old:=Sample{Name:"node1",RxBytes:1000,TxBytes:2000,Uptime:1000,Timestamp:time.Now().UTC().Add(-20*time.Second)}
 if err:=recordSample(db,old);err!=nil{t.Fatal(err)}
 next:=old;next.Timestamp=old.Timestamp.Add(10*time.Second);next.RxBytes=2400;next.TxBytes=2800;next.Uptime=1010
 if err:=recordTrafficIncrement(db,next);err!=nil{t.Fatal(err)}
 if err:=recordSample(db,next);err!=nil{t.Fatal(err)}
 reboot:=next;reboot.Timestamp=next.Timestamp.Add(10*time.Second);reboot.RxBytes=100;reboot.TxBytes=50;reboot.Uptime=3
 if err:=recordTrafficIncrement(db,reboot);err!=nil{t.Fatal(err)}
 settings:=NodeLimits{Name:"node1",Timezone:"Asia/Shanghai",QuotaGB:2,ExpiresOn:"2027-03-01"}
 if err:=setNodeLimits(context.Background(),db,settings);err!=nil{t.Fatal(err)}
 periods,err:=periodTraffic(context.Background(),db);if err!=nil{t.Fatal(err)}
 p:=periods["node1"]
 if !p.HasSamples || p.MonthRX!=1400 || p.MonthTX!=800 {t.Fatalf("unexpected traffic: %+v",p)}
 if p.Timezone!="Asia/Shanghai"||p.QuotaGB!=2||p.ExpiresOn!="2027-03-01"{t.Fatalf("metadata lost: %+v",p)}
}

func TestNodeLimitsValidation(t *testing.T) {
 for _,tz:=range []string{"UTC","Asia/Shanghai","Europe/Amsterdam","Asia/Kathmandu"}{
  if !validNodeLimits(NodeLimits{Name:"hk01",Timezone:tz,QuotaGB:500,ExpiresOn:"2027-01-10"}){t.Fatalf("timezone rejected: %s",tz)}
 }
 for _,v:=range []NodeLimits{
  {Name:"a",Timezone:"../../etc/passwd"},
  {Name:"a",Timezone:"Invalid/Whatever"},
  {Name:"a",Timezone:"UTC",QuotaGB:-1},
  {Name:"a",Timezone:"UTC",ExpiresOn:"2027-02-31"},
 }{
  if validNodeLimits(v){t.Fatalf("invalid config accepted: %+v",v)}
 }
}

func TestQuotaAndExpiryAlertStages(t *testing.T) {
 p:=PeriodTraffic{QuotaGB:1,HasSamples:true}
 for _,tc:=range []struct{used uint64; want string}{{700000000,""},{800000000,"quota_80"},{920000000,"quota_90"},{1000000000,"quota_100"}}{
  p.MonthRX=tc.used
  if got:=quotaKind(p);got!=tc.want{t.Fatalf("got %s want %s",got,tc.want)}
 }
 now:=time.Date(2026,10,10,12,0,0,0,time.UTC)
 for _,tc:=range []struct{expiry,kind string}{{"2026-11-30",""},{"2026-11-05","expiry_30"},{"2026-10-24","expiry_15"},{"2026-10-15","expiry_7"},{"2026-10-09","expiry_overdue"}}{
  if got:=expirationKind(tc.expiry,"UTC",now);got!=tc.kind{t.Fatalf("expiry %s got %s want %s",tc.expiry,got,tc.kind)}
 }
}
