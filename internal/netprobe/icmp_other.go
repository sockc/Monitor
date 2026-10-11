//go:build !windows && !linux

package netprobe

import "context"

func measureICMP(context.Context,string)*float64{return nil}
