package main

import (
 "io"
 "net"
 "net/http"
 "strings"
 "testing"
 "time"
)

type mockTransport func(*http.Request)(*http.Response,error)
func (fn mockTransport) RoundTrip(r *http.Request)(*http.Response,error){return fn(r)}
func TestAgentGeoUsesHTTPSAndValidatesResult(t *testing.T){
 client:=&http.Client{Timeout:time.Second,Transport:mockTransport(func(r *http.Request)(*http.Response,error){
  if r.URL.Scheme!="https"||r.URL.Host!="ipwho.is"{t.Errorf("unexpected geo endpoint %s",r.URL.String())}
  return &http.Response{StatusCode:200,Body:io.NopCloser(strings.NewReader(`{"success":true,"ip":"8.8.8.8","country":"United States","city":"Mountain View"}`)),Header:make(http.Header)},nil
 })}
 got,err:=lookupAgentGeo(client);if err!=nil{t.Fatal(err)}
 if got.IP!="8.8.8.8"||got.Location!="United States · Mountain View"{t.Fatalf("bad auto geo: %+v",got)}
}
func TestGeoInvalidAddress(t *testing.T){
 for _,ip:=range []string{"","localhost","127.0.0.1","192.168.2.10","10.0.1.2","::1","fe80::1","0.0.0.0"}{
  if validPublicIP(ip){t.Fatalf("private IP accepted: %q",ip)}
 }
 if !validPublicIP("8.8.4.4")||!validPublicIP("2001:4860:4860::8888"){t.Fatal("public IPv4 or IPv6 rejected")}
}
func TestAgentGeoProviderFailure(t *testing.T){
 client:=&http.Client{Transport:mockTransport(func(r *http.Request)(*http.Response,error){
  return &http.Response{StatusCode:429,Body:io.NopCloser(strings.NewReader("rate limited")),Header:make(http.Header)},nil
 })}
 if _,err:=lookupAgentGeo(client);err==nil{t.Fatal("expected rate-limit error")}
}
var _=net.IPv4
