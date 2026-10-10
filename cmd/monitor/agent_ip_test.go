package main

import (
 "io"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestOutboundIPFamilies(t *testing.T){
 cases:=[]struct{family,body,want string}{
  {"4","8.8.4.4\n","8.8.4.4"},
  {"6","2001:4860:4860::8888\n","2001:4860:4860::8888"},
 }
 for _,tc:=range cases{
  t.Run(tc.family,func(t *testing.T){
   client:=&http.Client{Transport:mockTransport(func(req *http.Request)(*http.Response,error){
    if req.URL.Scheme!="https"||req.URL.Host!="api"+tc.family+".ipify.org"{t.Fatalf("unexpected URL: %s",req.URL)}
    return &http.Response{StatusCode:200,Body:io.NopCloser(strings.NewReader(tc.body)),Header:make(http.Header)},nil
   })}
   got,err:=lookupIPFamily(client,tc.family)
   if err!=nil||got!=tc.want{t.Fatalf("family %s got %q, err=%v",tc.family,got,err)}
  })
 }
}
func TestOutboundIPFailuresAreUnknown(t *testing.T){
 for _,tc:=range []struct{family,ip string}{{"4","2001:4860:4860::8888"},{"6","8.8.8.8"},{"4","127.0.0.1"},{"6","fe80::1"},{"4","invalid"}}{
  client:=&http.Client{Transport:mockTransport(func(req *http.Request)(*http.Response,error){
   return &http.Response{StatusCode:200,Body:io.NopCloser(strings.NewReader(tc.ip)),Header:make(http.Header)},nil
  })}
  if _,err:=lookupIPFamily(client,tc.family);err==nil{t.Fatalf("should reject family %s, IP %s",tc.family,tc.ip)}
 }
}
func TestSelfHostedFlags(t *testing.T){
 for _,code:=range []string{"gb","us","nl","hk","de"}{
  req:=httptest.NewRequest(http.MethodGet,"/static/flags/"+code+".svg",nil)
  res:=httptest.NewRecorder()
  staticHandler().ServeHTTP(res,req)
  if res.Code!=200||!strings.Contains(res.Body.String(),"<svg"){t.Fatalf("missing %s flag: HTTP %d",code,res.Code)}
  if res.Header().Get("Content-Type")!="image/svg+xml; charset=utf-8"{t.Fatalf("invalid flag content-type: %s",code)}
 }
}
