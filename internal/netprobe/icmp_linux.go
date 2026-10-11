//go:build linux

package netprobe

import (
 "context"
 "os/exec"
 "time"
)

func measureICMP(ctx context.Context,server string)*float64 {
 host:=pingHost(server)
 if host==""{return nil}
 if _,err:=exec.LookPath("ping");err!=nil{return nil}
 sub,cancel:=context.WithTimeout(ctx,6*time.Second);defer cancel()
 // argv is never interpolated through a shell.
 output,_:=exec.CommandContext(sub,"ping","-c","3","-W","1","-q",host).CombinedOutput()
 return parseICMPLoss(string(output))
}
