package main

import (
 "testing"
 "time"
)

func TestBootTimeFromSample(t *testing.T) {
 reported:=time.Date(2026,10,10,12,30,42,0,time.UTC)
 sample:=Sample{Timestamp:reported,Uptime:2*24*3600+3*3600+19}
 got:=bootTimeFromSample(sample)
 want:="2026-10-08T09:30:23Z"
 if got!=want { t.Fatalf("inferred boot = %q, expected %q",got,want) }
 if bootTimeFromSample(Sample{Timestamp:reported})!="" { t.Fatal("missing uptime must not invent boot time") }
 if bootTimeFromSample(Sample{Uptime:3600})!="" { t.Fatal("missing timestamp must not invent boot time") }
 if bootTimeFromSample(Sample{Timestamp:reported,Uptime:150*366*24*3600})!="" { t.Fatal("unrealistically old uptime must not produce boot time") }
}

func TestBootTimeUsesAgentClock(t *testing.T){
 // A remote server may have a different timezone than the dashboard viewer.
 // Expose UTC RFC3339; JavaScript formats it in the viewer's own timezone.
 zone:=time.FixedZone("UTC+8",8*3600)
 sample:=Sample{Timestamp:time.Date(2026,10,10,20,0,0,0,zone),Uptime:3600}
 if got:=bootTimeFromSample(sample);got!="2026-10-10T11:00:00Z"{t.Fatal(got)}
}
