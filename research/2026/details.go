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
   req.URL.Host = "play.google.com"
   req.URL.Path = "/store/apps/details"
   value := url.Values{}
   value["id"] = []string{"com.wbd.stream"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      panic(err)
   }
   if strings.Contains(data.String(), "100,000,000") {
      log.Print("pass")
   } else {
      log.Print("fail")
   }
}
