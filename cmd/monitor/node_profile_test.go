package main

import (
 "context"
 "database/sql"
 "path/filepath"
 "testing"
)

func TestProfileFieldsValidation(t *testing.T) {
 p:=defaultNodeProfile()
 if !validNodeProfile(p){t.Fatal("empty valid profile rejected")}
 p.Provider="CloudCone"
 p.CountryCode="US"
 p.PriceValue=14.5
 p.PriceCurrency="USD"
 p.BillingCycle="year"
 p.PeriodStart="2026-01-01"
 p.PortMbps=1000
 p.HasIPv4=1
 p.HasIPv6=0
 if !validNodeProfile(p){t.Fatalf("valid full profile rejected: %+v",p)}
 for _,bad:=range []NodeProfile{
  func()NodeProfile{x:=p;x.CountryCode="USA";return x}(),
  func()NodeProfile{x:=p;x.PriceValue= -1;return x}(),
  func()NodeProfile{x:=p;x.PortMbps= -1;return x}(),
  func()NodeProfile{x:=p;x.HasIPv6=2;return x}(),
  func()NodeProfile{x:=p;x.PriceCurrency="BAD";return x}(),
  func()NodeProfile{x:=p;x.PeriodStart="2026-02-30";return x}(),
 }{if validNodeProfile(bad){t.Fatalf("invalid profile accepted: %+v",bad)}}
}

func TestProfileMigrationRenameAndDeletion(t *testing.T) {
 ctx:=context.Background()
 db,err:=openDB(filepath.Join(t.TempDir(),"old-monitor.db"));if err!=nil{t.Fatal(err)};defer db.Close()
 _,err=db.Exec("INSERT INTO node_profile(name,provider,country_code,price_value,price_currency,billing_cycle,period_start,port_mbps,has_ipv4,has_ipv6) VALUES('us1','CloudCone','US',14.50,'USD','year','2026-01-01',1000,1,1)")
 if err!=nil{t.Fatal(err)}
 _,err=db.Exec("INSERT INTO node_geo(name,public_ip,auto_location,country_code,updated_at) VALUES('us1','8.8.8.8','United States','US',12345)")
 if err!=nil{t.Fatal(err)}
 profiles,err:=readNodeProfiles(ctx,db);if err!=nil{t.Fatal(err)}
 if p:=profiles["us1"];p.PortMbps!=1000||p.PriceValue!=14.5||p.PeriodStart!="2026-01-01"||p.HasIPv4!=1||p.HasIPv6!=1{t.Fatalf("incorrect profile: %+v",p)}
 if err=renameNode(ctx,db,"us1","us2");err!=nil{t.Fatal(err)}
 profiles,err=readNodeProfiles(ctx,db);if err!=nil{t.Fatal(err)}
 if _,ok:=profiles["us1"];ok{t.Fatal("old profile remains")}
 if profiles["us2"].Provider!="CloudCone"{t.Fatal("profile missing after rename")}
 var country string
 if err=db.QueryRow("SELECT country_code FROM node_geo WHERE name='us2'").Scan(&country);err!=nil||country!="US"{t.Fatalf("geo migration: %q, %v",country,err)}
 if err=removeNode(ctx,db,"us2");err!=nil{t.Fatal(err)}
 if err=db.QueryRow("SELECT country_code FROM node_geo WHERE name='us2'").Scan(&country);err!=sql.ErrNoRows{t.Fatalf("geo not removed: %v",err)}
 profiles,err=readNodeProfiles(ctx,db);if err!=nil{t.Fatal(err)}
 if len(profiles)!=0{t.Fatalf("orphan profiles: %+v",profiles)}
}
