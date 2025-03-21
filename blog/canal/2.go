package main

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "tvapi-hlm2.solocoo.tv"
   req.URL.Path = "/v1/session"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(body)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

var body = strings.NewReader(`
{
   "brand": "m7cp",
   "deviceSerial": "w1f8a8fb0-05fb-11f0-b5da-f32d1cd5ddf7",
   "deviceType": "PC",
   "ssoToken": "eyJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2V5IjoibTciLCJhbGciOiJkaXIifQ..pU54fOS8W-V9GhMJHdpVyw.uxmy6DiLsvodKRl3CMRmVNEEssfNZd9opOyDcOuM_MQONhDwAnYoNwt5MZHMXwelsIZEcjVDxiBCQ1Yy1QSLkjyFDqLkKC6kuid_2WSmIYMkuPkAaXdNrL8SkpHScyb6aBvzCL2qGD3uO6ElWYfGJ3cJCCrHbguMkKhbiO3EPB2Ng16kEBxmqAGnCHb2O04M5q3JxOwzxQLW8G1chGiOdGyG3nrlc_-PKWzdU8JJI8PUvTmFh8AzM0D6siCXbCaKRRP8OT3ek0JI9G1Rlc581TgMtNWrwuAPP2Vi1sF-WlCcCGGM3R0mUNCKkCJZ_04m52C8IKn9kF3Ka6oGWloeF7IgcPqlv5lPcqg9GXJ3RtbntFYIRzB_5Pj61ADAmIGUNI84U7qu3tNQEQy4GGBKmkBYRwmkkMulbvJMv33TENy5TEVpgMsQDj-lCubi2Cex.WsfESmeAt9pKxc5ifFshBg"
}
`)
