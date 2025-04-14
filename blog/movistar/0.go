package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header["Accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
   req.Header["Accept-Language"] = []string{"es-ES,es;q=0.9"}
   req.Header["Connection"] = []string{"keep-alive"}
   req.Header["Content-Length"] = []string{"0"}
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Host"] = []string{"auth.dof6.com"}
   req.Header["Origin"] = []string{"https://ver.movistarplus.es"}
   req.Header["Referer"] = []string{"https://ver.movistarplus.es/"}
   req.Header["User-Agent"] = []string{"Dalvik/2.1.0 (Linux; U; Android 12; 22126RN91Y Build/SP1A.210812.016)"}
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "auth.dof6.com"
   value := url.Values{}
   value["qspVersion"] = []string{"ssp"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   
   req.Header["Authorization"] = []string{"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI3VTdlN3Y4QjhTOGg4bzlBIiwiYWNjb3VudE51bWJlciI6IjAwUVNwMDAwMDA5TTlnek1BQy1MIiwicm9sZSI6InVzZXIiLCJhcHIiOiJ3ZWJkYiIsImlzcyI6Imh0dHA6Ly93d3cubW92aXN0YXJwbHVzLmVzIiwiYXVkIjoiNDE0ZTE5MjdhMzg4NGY2OGFiYzc5ZjcyODM4MzdmZDEiLCJleHAiOjE3NDUyNzYwMzQsIm5iZiI6MTc0NDQxMjAzNH0.HdHCbmJaU-a5ASyrrygJRQF9RHj9cdN9eyx9A341H3s"}
   req.Header["Cookie"] = []string{"080a03=zrMaQO9lCknTAL21WupGekZHsFMVkpMbpalnCLuvyeE+6ZByZAKP9LNxFA65tN4OmDEhNkdsha3YIBLcBmGyAMRcnJwwGu0hxSvzP7Rfs7Z2UR9eVSniweenaAeHlGxpFQ5LvCydEijXaLsNUWAuzX+HGIBi7gUFdb1SMiV1f7ybhdUS"}
   req.URL.Path = "/movistarplus/amazon.tv/accounts/00QSp000009M9gzMAC-L/devices/"
   
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}
