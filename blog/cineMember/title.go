package main

import (
   "io"
   "net/http"
   "net/url"
   "strings"
   "log"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "www.cinemember.nl"
   req.URL.Path = "/nl/title/468545/american-hustle"
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      panic(err)
   }
   if strings.Contains(string(data), "906945") {
      log.Print("pass")
   } else {
      log.Print("fail")
   }
}
