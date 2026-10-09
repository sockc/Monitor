//go:build !linux

package main

import ("os";"runtime";"time")
type cpuTimes struct{}
func readCPU()cpuTimes{return cpuTimes{}}
func collect(name string,prev cpuTimes)(Sample,cpuTimes){h,_:=os.Hostname();return Sample{Name:name,Hostname:h,OS:runtime.GOOS,Arch:runtime.GOARCH,Timestamp:time.Now().UTC()},prev}
