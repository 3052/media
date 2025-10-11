package main

import (
   "net/http"
   "net/url"
   "fmt"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "www.cinemember.nl"
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", resp)
}
