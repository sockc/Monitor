//go:build windows

package netprobe

import (
 "context"
 "os"
 "os/exec"
 "path/filepath"
 "time"
)

func measureICMP(ctx context.Context,server string)*float64 {
 host:=pingHost(server)
 if host==""{return nil}
 exe:="ping.exe"
 if _,err:=exec.LookPath(exe);err!=nil {
  root:=os.Getenv("SystemRoot")
  if root==""{return nil}
  exe=filepath.Join(root,"System32","ping.exe")
 }
 sub,cancel:=context.WithTimeout(ctx,7*time.Second);defer cancel()
 output,_:=exec.CommandContext(sub,exe,"-n","3","-w","1000",host).CombinedOutput()
 return parseICMPLoss(string(output))
}
