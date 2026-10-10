//go:build !windows

package main

import "time"

func runAgentPlatform(server,name,token string,interval time.Duration){
 runAgent(server,name,token,interval)
}
