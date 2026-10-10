package main

import (
 "fmt"
 "net/http"
 "strings"
)

// SVG designs from lipis/flag-icons (MIT License), https://github.com/lipis/flag-icons
// Icons are served locally so Windows browsers without colored flag emojis work.
var flagSVGs=map[string]string{
"GB":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-gb\" viewBox=\"0 0 640 480\">\n  <path fill=\"#012169\" d=\"M0 0h640v480H0z\"/>\n  <path fill=\"#FFF\" d=\"m75 0 244 181L562 0h78v62L400 241l240 178v61h-80L320 301 81 480H0v-60l239-178L0 64V0z\"/>\n  <path fill=\"#C8102E\" d=\"m424 281 216 159v40L369 281zm-184 20 6 35L54 480H0zM640 0v3L391 191l2-44L590 0zM0 0l239 176h-60L0 42z\"/>\n  <path fill=\"#FFF\" d=\"M241 0v480h160V0zM0 160v160h640V160z\"/>\n  <path fill=\"#C8102E\" d=\"M0 193v96h640v-96zM273 0v480h96V0z\"/>\n</svg>\n",
"US":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-us\" viewBox=\"0 0 640 480\">\n  <path fill=\"#bd3d44\" d=\"M0 0h640v480H0\"/>\n  <path stroke=\"#fff\" stroke-width=\"37\" d=\"M0 55.3h640M0 129h640M0 203h640M0 277h640M0 351h640M0 425h640\"/>\n  <path fill=\"#192f5d\" d=\"M0 0h364.8v258.5H0\"/>\n  <marker id=\"us-a\" markerHeight=\"30\" markerWidth=\"30\">\n    <path fill=\"#fff\" d=\"m14 0 9 27L0 10h28L5 27z\"/>\n  </marker>\n  <path fill=\"none\" marker-mid=\"url(#us-a)\" d=\"m0 0 16 11h61 61 61 61 60L47 37h61 61 60 61L16 63h61 61 61 61 60L47 89h61 61 60 61L16 115h61 61 61 61 60L47 141h61 61 60 61L16 166h61 61 61 61 60L47 192h61 61 60 61L16 218h61 61 61 61 60z\"/>\n</svg>\n",
"NL":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-nl\" viewBox=\"0 0 640 480\">\n  <path fill=\"#ae1c28\" d=\"M0 0h640v160H0z\"/>\n  <path fill=\"#fff\" d=\"M0 160h640v160H0z\"/>\n  <path fill=\"#21468b\" d=\"M0 320h640v160H0z\"/>\n</svg>\n",
"HK":"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" id=\"flag-icons-hk\" viewBox=\"0 0 640 480\">\n  <path fill=\"#EC1B2E\" d=\"M0 0h640v480H0\"/>\n  <path id=\"hk-a\" fill=\"#fff\" d=\"M346.3 103.1C267 98 230.6 201.9 305.6 240.3c-26-22.4-20.6-55.3-10.1-72.4l1.9 1.1c-13.8 23.5-11.2 52.7 11.1 71-12.7-12.3-9.5-39 12.1-48.9s23.6-39.3 16.4-49.1q-14.7-25.6 9.3-38.9M307.9 164l-4.7 7.4-1.8-8.6-8.6-2.3 7.8-4.3-.6-8.9 6.5 6.1 8.3-3.3-3.7 8.1 5.6 6.8z\"/>\n  <use xlink:href=\"#hk-a\" transform=\"rotate(72 312.5 243.5)\"/>\n  <use xlink:href=\"#hk-a\" transform=\"rotate(144 312.5 243.5)\"/>\n  <use xlink:href=\"#hk-a\" transform=\"rotate(216 312.5 243.5)\"/>\n  <use xlink:href=\"#hk-a\" transform=\"rotate(288 312.5 243.5)\"/>\n</svg>\n",
"CN":"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" id=\"flag-icons-cn\" viewBox=\"0 0 640 480\">\n  <defs>\n    <path id=\"cn-a\" fill=\"#ff0\" d=\"M-.6.8 0-1 .6.8-1-.3h2z\"/>\n  </defs>\n  <path fill=\"#ee1c25\" d=\"M0 0h640v480H0z\"/>\n  <use xlink:href=\"#cn-a\" width=\"30\" height=\"20\" transform=\"matrix(71.9991 0 0 72 120 120)\"/>\n  <use xlink:href=\"#cn-a\" width=\"30\" height=\"20\" transform=\"matrix(-12.33562 -20.5871 20.58684 -12.33577 240.3 48)\"/>\n  <use xlink:href=\"#cn-a\" width=\"30\" height=\"20\" transform=\"matrix(-3.38573 -23.75998 23.75968 -3.38578 288 95.8)\"/>\n  <use xlink:href=\"#cn-a\" width=\"30\" height=\"20\" transform=\"matrix(6.5991 -23.0749 23.0746 6.59919 288 168)\"/>\n  <use xlink:href=\"#cn-a\" width=\"30\" height=\"20\" transform=\"matrix(14.9991 -18.73557 18.73533 14.99929 240 216)\"/>\n</svg>\n",
"JP":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-jp\" viewBox=\"0 0 640 480\">\n  <defs>\n    <clipPath id=\"jp-a\">\n      <path fill-opacity=\".7\" d=\"M-88 32h640v480H-88z\"/>\n    </clipPath>\n  </defs>\n  <g fill-rule=\"evenodd\" stroke-width=\"1pt\" clip-path=\"url(#jp-a)\" transform=\"translate(88 -32)\">\n    <path fill=\"#fff\" d=\"M-128 32h720v480h-720z\"/>\n    <circle cx=\"523.1\" cy=\"344.1\" r=\"194.9\" fill=\"#bc002d\" transform=\"translate(-168.4 8.6)scale(.76554)\"/>\n  </g>\n</svg>\n",
"DE":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-de\" viewBox=\"0 0 640 480\">\n  <path fill=\"#fc0\" d=\"M0 320h640v160H0z\"/>\n  <path fill=\"#000001\" d=\"M0 0h640v160H0z\"/>\n  <path fill=\"red\" d=\"M0 160h640v160H0z\"/>\n</svg>\n",
"FR":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-fr\" viewBox=\"0 0 640 480\">\n  <path fill=\"#000091\" d=\"M0 0h213.3v480H0z\"/>\n  <path fill=\"#fff\" d=\"M213.3 0h213.4v480H213.3z\"/>\n  <path fill=\"#e1000f\" d=\"M426.7 0H640v480H426.7z\"/>\n</svg>\n",
"SG":"<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"flag-icons-sg\" viewBox=\"0 0 640 480\">\n  <defs>\n    <clipPath id=\"sg-a\">\n      <path fill-opacity=\".7\" d=\"M0 0h640v480H0z\"/>\n    </clipPath>\n  </defs>\n  <g fill-rule=\"evenodd\" clip-path=\"url(#sg-a)\">\n    <path fill=\"#fff\" d=\"M-20 0h720v480H-20z\"/>\n    <path fill=\"#df0000\" d=\"M-20 0h720v240H-20z\"/>\n    <path fill=\"#fff\" d=\"M146 40.2a84.4 84.4 0 0 0 .8 165.2 86 86 0 0 1-106.6-59 86 86 0 0 1 59-106c16-4.6 30.8-4.7 46.9-.2z\"/>\n    <path fill=\"#fff\" d=\"m133 110 4.9 15-13-9.2-12.8 9.4 4.7-15.2-12.8-9.3 15.9-.2 5-15 5 15h15.8zm17.5 52 5 15.1-13-9.2-12.9 9.3 4.8-15.1-12.8-9.4 15.9-.1 4.9-15.1 5 15h16zm58.5-.4 4.9 15.2-13-9.3-12.8 9.3 4.7-15.1-12.8-9.3 15.9-.2 5-15 5 15h15.8zm17.4-51.6 4.9 15.1-13-9.2-12.8 9.3 4.8-15.1-12.9-9.4 16-.1 4.8-15.1 5 15h16zm-46.3-34.3 5 15.2-13-9.3-12.9 9.4 4.8-15.2-12.8-9.4 15.8-.1 5-15.1 5 15h16z\"/>\n  </g>\n</svg>\n",
"KR":"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" id=\"flag-icons-kr\" viewBox=\"0 0 640 480\">\n  <defs>\n    <clipPath id=\"kr-a\">\n      <path fill-opacity=\".7\" d=\"M-95.8-.4h682.7v512H-95.8z\"/>\n    </clipPath>\n  </defs>\n  <g fill-rule=\"evenodd\" clip-path=\"url(#kr-a)\" transform=\"translate(89.8 .4)scale(.9375)\">\n    <path fill=\"#fff\" d=\"M-95.8-.4H587v512H-95.8Z\"/>\n    <g transform=\"rotate(-56.3 361.6 -101.3)scale(10.66667)\">\n      <g id=\"kr-c\">\n        <path id=\"kr-b\" fill=\"#000001\" d=\"M-6-26H6v2H-6Zm0 3H6v2H-6Zm0 3H6v2H-6Z\"/>\n        <use xlink:href=\"#kr-b\" width=\"100%\" height=\"100%\" y=\"44\"/>\n      </g>\n      <path stroke=\"#fff\" d=\"M0 17v10\"/>\n      <path fill=\"#cd2e3a\" d=\"M0-12a12 12 0 0 1 0 24Z\"/>\n      <path fill=\"#0047a0\" d=\"M0-12a12 12 0 0 0 0 24A6 6 0 0 0 0 0Z\"/>\n      <circle cy=\"-6\" r=\"6\" fill=\"#cd2e3a\"/>\n    </g>\n    <g transform=\"rotate(-123.7 191.2 62.2)scale(10.66667)\">\n      <use xlink:href=\"#kr-c\" width=\"100%\" height=\"100%\"/>\n      <path stroke=\"#fff\" d=\"M0-23.5v3M0 17v3.5m0 3v3\"/>\n    </g>\n  </g>\n</svg>\n",
}

func serveFlagSVG(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet && r.Method!=http.MethodHead{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
 code:=strings.TrimPrefix(r.URL.Path,"/static/flags/")
 if len(code)!=6 || !strings.HasSuffix(code,".svg"){http.NotFound(w,r);return}
 code=strings.ToUpper(strings.TrimSuffix(code,".svg"))
 for _,c:=range code{if c<'A'||c>'Z'{http.NotFound(w,r);return}}
 data,ok:=flagSVGs[code]
 if !ok{
   // A readable code marker is safer than rendering a false national flag.
   data=fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 480"><rect width="640" height="480" fill="#e5edfa"/><rect x="12" y="12" width="616" height="456" rx="30" fill="none" stroke="#8194b5" stroke-width="20"/><text x="320" y="310" font-family="Arial,sans-serif" font-size="210" text-anchor="middle" font-weight="bold" fill="#475e82">%s</text></svg>`,code)
 }
 w.Header().Set("Content-Type","image/svg+xml; charset=utf-8")
 w.Header().Set("X-Content-Type-Options","nosniff")
 w.Header().Set("Cache-Control","public, max-age=86400")
 w.Write([]byte(data))
}
