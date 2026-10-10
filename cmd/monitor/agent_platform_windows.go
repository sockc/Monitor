//go:build windows

package main

import (
 "context"
 "log"
 "time"

 "golang.org/x/sys/windows/svc"
)

type monitorAgentService struct {server,name,token string;interval time.Duration}
func (m *monitorAgentService) Execute(_ []string,req <-chan svc.ChangeRequest,status chan<- svc.Status)(bool,uint32){
 status<-svc.Status{State:svc.StartPending}
 ctx,cancel:=context.WithCancel(context.Background())
 done:=make(chan struct{})
 go func(){defer close(done);runAgentContext(ctx,m.server,m.name,m.token,m.interval)}()
 status<-svc.Status{State:svc.Running,Accepts:svc.AcceptStop|svc.AcceptShutdown}
 for change:=range req {
  switch change.Cmd{
  case svc.Interrogate:status<-change.CurrentStatus
  case svc.Stop,svc.Shutdown:
   status<-svc.Status{State:svc.StopPending}
   cancel()
   select{case <-done:case <-time.After(12*time.Second):}
   return false,0
  }
 }
 cancel()
 return false,0
}
func runAgentPlatform(server,name,token string,interval time.Duration){
 isService,e:=svc.IsWindowsService()
 if e!=nil{log.Fatalf("detect Windows service: %v",e)}
 if isService{
  if e=svc.Run("MonitorAgent",&monitorAgentService{server:server,name:name,token:token,interval:interval});e!=nil{log.Fatalf("Windows service failed: %v",e)}
  return
 }
 runAgent(server,name,token,interval)
}
