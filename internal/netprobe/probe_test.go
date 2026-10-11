package netprobe

import (
 "context"
 "errors"
 "net"
 "testing"
)

func TestICMPLossParser(t *testing.T) {
 for _,tc:=range []struct{text string;expected float64;known bool}{
  {"3 packets transmitted, 2 received, 33.3% packet loss",33.3,true},
  {"Packets: Sent = 3, Received = 3, Lost = 0 (0% loss)",0,true},
  {"Packets: Sent = 3, Received = 0, Lost = 3 (100% loss)",100,true},
  {"ping command unavailable",0,false},
 }{
  got:=parseICMPLoss(tc.text)
  if (got!=nil)!=tc.known||(got!=nil&&*got!=tc.expected){t.Errorf("unexpected parse for %q: %v",tc.text,got)}
 }
}
func TestEndpoint(t *testing.T) {
 for _,tt:=range []struct{server,want string;ok bool}{
  {"https://monitor.example.com","monitor.example.com:443",true},
  {"https://[2001:db8::1]:8443","[2001:db8::1]:8443",true},
  {"https://monitor.example.com:9443/app","monitor.example.com:9443",true},
  {"http://monitor.example.com","",false},
  {"https://user:secret@example.com","",false},
  {"https://bad.example:abc","",false},
  {"http://127.0.0.1:8080","127.0.0.1:8080",true},
 } {
  got,err:=Endpoint(tt.server)
  if (err==nil)!=tt.ok||got!=tt.want {t.Errorf("%q => %q, %v",tt.server,got,err)}
 }
}
func TestMeasureFailureNotICMPLoss(t *testing.T) {
 fail:=func(context.Context,string,string)(net.Conn,error){return nil,errors.New("refused")}
 got:=Measure(context.Background(),"localhost:443",3,fail)
 if got.FailurePercent==nil||*got.FailurePercent!=100||got.RTTMS!=nil||got.CheckedAt==0{t.Fatalf("%+v",got)}
}
func TestMeasureSuccess(t *testing.T){
 l,e:=net.Listen("tcp","127.0.0.1:0");if e!=nil{t.Fatal(e)}
 defer l.Close()
 got:=Measure(context.Background(),l.Addr().String(),3,(&net.Dialer{}).DialContext)
 if got.RTTMS==nil||got.FailurePercent==nil||*got.FailurePercent!=0{t.Fatalf("%+v",got)}
}
