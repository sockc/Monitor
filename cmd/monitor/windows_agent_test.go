package main

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestAgentConfigFile(t *testing.T){
 for _,key:=range []string{"MONITOR_SERVER","MONITOR_NODE_NAME","MONITOR_AGENT_TOKEN"}{t.Setenv(key,"")}
 path:=filepath.Join(t.TempDir(),"agent.env")
 content:="MONITOR_SERVER=https://monitor.example.com\r\nMONITOR_NODE_NAME=n-test\r\nMONITOR_AGENT_TOKEN="+strings.Repeat("a",64)+"\r\n"
 if err:=os.WriteFile(path,[]byte(content),0600);err!=nil{t.Fatal(err)}
 if err:=loadAgentConfig(path);err!=nil{t.Fatal(err)}
 if got:=os.Getenv("MONITOR_SERVER");got!="https://monitor.example.com"{t.Fatalf("server: %q",got)}
 if got:=os.Getenv("MONITOR_NODE_NAME");got!="n-test"{t.Fatalf("node: %q",got)}
 if got:=os.Getenv("MONITOR_AGENT_TOKEN");got!=strings.Repeat("a",64){t.Fatal("token mismatch")}
}

func TestAgentConfigRejectsUnexpectedInput(t *testing.T){
 for _,body:=range []string{
  "UNKNOWN_ENV=evil\n",
  "MONITOR_AGENT_TOKEN=\n",
  "MALFORMED\n",
  strings.Repeat("a",8200),
 }{
  path:=filepath.Join(t.TempDir(),"agent.env")
  if err:=os.WriteFile(path,[]byte(body),0600);err!=nil{t.Fatal(err)}
  if err:=loadAgentConfig(path);err==nil{t.Fatalf("invalid config accepted (length=%d)",len(body))}
 }
}

func TestWindowsInstallPlatformInDashboard(t *testing.T){
 for _,piece:=range []string{
  `id="install-target"`,
  `value="windows"`,
  `value="linux"`,
 }{
  if !strings.Contains(dashboardHTML,piece){t.Errorf("missing dashboard: %s",piece)}
 }
 for _,piece:=range []string{
  "function generatedInstallCommand(url,id,token,target)",
  "install-windows.ps1",
  "-Action Install",
  "generatedInstallCommand(url,data.name,data.token,$('install-target').value)",
 }{
  if !strings.Contains(appJS,piece){t.Errorf("missing Windows command generator: %s",piece)}
 }
}
