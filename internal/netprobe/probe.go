package netprobe

import (
 "context"
 "fmt"
 "math"
 "net"
 "net/url"
 "strconv"
 "sync"
 "time"
)

// Result describes TCP connect quality to the configured Monitor server.
// FailurePercent is a TCP-connect failure ratio, NOT ICMP packet loss.
type Result struct {
 RTTMS *float64 `json:"network_rtt_ms,omitempty"`
 FailurePercent *float64 `json:"network_failure_percent,omitempty"`
 ICMPLossPercent *float64 `json:"network_icmp_loss_percent,omitempty"`
 CheckedAt int64 `json:"network_checked_at,omitempty"`
}

func Endpoint(server string) (string,error) {
 u,e:=url.Parse(server)
 if e!=nil||u.Hostname()==""||u.User!=nil{return "",fmt.Errorf("invalid Monitor server URL")}
 if u.Scheme!="https"&&!(u.Scheme=="http"&&u.Hostname()=="127.0.0.1"){return "",fmt.Errorf("HTTPS required")}
 port:=u.Port();if port==""{if u.Scheme=="https"{port="443"}else{port="80"}}
 n,e:=strconv.Atoi(port);if e!=nil||n<1||n>65535{return "",fmt.Errorf("invalid server port")}
 return net.JoinHostPort(u.Hostname(),port),nil
}

func Measure(ctx context.Context,address string, attempts int, dial func(context.Context,string,string)(net.Conn,error)) Result {
 if attempts<1 {attempts=1};if attempts>5{attempts=5}
 out:=Result{CheckedAt:time.Now().Unix()}
 successes:=0;var total float64
 for i:=0;i<attempts;i++ {
  if ctx.Err()!=nil{break}
  sub,cancel:=context.WithTimeout(ctx,1200*time.Millisecond)
  start:=time.Now()
  conn,err:=dial(sub,"tcp",address)
  elapsed:=time.Since(start).Seconds()*1000
  cancel()
  if err==nil{
   conn.Close()
   total+=elapsed
   successes++
  }
 }
 failure:=float64(attempts-successes)*100/float64(attempts)
 out.FailurePercent=&failure
 if successes>0 {rtt:=math.Round(total/float64(successes)*10)/10;out.RTTMS=&rtt}
 return out
}

type Sampler struct{
 mu sync.RWMutex
 result Result
}

func (s *Sampler) Start(ctx context.Context,server string) {
 address,e:=Endpoint(server)
 if e!=nil{return}
 go func(){
  dialer:=&net.Dialer{}
  ticker:=time.NewTicker(60*time.Second);defer ticker.Stop()
  for {
   result:=Measure(ctx,address,3,dialer.DialContext)
   result.ICMPLossPercent=measureICMP(ctx,server)
   s.mu.Lock();s.result=result;s.mu.Unlock()
   select {case <-ctx.Done():return;case <-ticker.C:}
  }
 }()
}

// Snapshot never returns stale quality as if it were current.
func (s *Sampler) Snapshot() Result {
 s.mu.RLock();p:=s.result;s.mu.RUnlock()
 if p.CheckedAt==0||time.Since(time.Unix(p.CheckedAt,0))>180*time.Second{return Result{}}
 return p
}
