package netprobe

import (
 "net/url"
 "regexp"
 "strconv"
 "strings"
)

var icmpLossPattern=regexp.MustCompile(`([0-9]+(?:[.,][0-9]+)?)\s*%`)

func parseICMPLoss(output string)*float64 {
 match:=icmpLossPattern.FindStringSubmatch(output)
 if len(match)!=2{return nil}
 v,err:=strconv.ParseFloat(strings.ReplaceAll(match[1],",","."),64)
 if err!=nil||v<0||v>100{return nil}
 return &v
}
func pingHost(server string)string{
 u,err:=url.Parse(server)
 if err!=nil{return ""}
 return u.Hostname()
}
