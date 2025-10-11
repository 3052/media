package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "www.cinemember.nl"
   req.URL.Path = "/elements/films/stream.php"
   req.URL.Scheme = "https"
   req.Header.Add("Cookie", "PHPSESSID=f2c0dcaf9b32eebe2fd0038cc58b776c")
   value := url.Values{}
   value["id"] = []string{"917398"}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}
