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
   req.Header["Accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
   req.Header["Accept-Language"] = []string{"es-ES,es;q=0.9"}
   req.Header["Connection"] = []string{"keep-alive"}
   req.Header["Content-Length"] = []string{"204"}
   req.ContentLength = 204
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Host"] = []string{"clientservices.dof6.com"}
   req.Header["Origin"] = []string{"https://ver.movistarplus.es"}
   req.Header["Referer"] = []string{"https://ver.movistarplus.es/"}
   req.Header["User-Agent"] = []string{"Dalvik/2.1.0 (Linux; U; Android 12; 22126RN91Y Build/SP1A.210812.016)"}
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "clientservices.dof6.com"
   value := url.Values{}
   value["qspVersion"] = []string{"ssp"}
   value["status"] = []string{"default"}
   value["version"] = []string{"8"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(body)
   req.Header["Authorization"] = []string{"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI3VTdlN3Y4QjhTOGg4bzlBIiwiYWNjb3VudE51bWJlciI6IjAwUVNwMDAwMDA5TTlnek1BQy1MIiwicm9sZSI6InVzZXIiLCJhcHIiOiJ3ZWJkYiIsImlzcyI6Imh0dHA6Ly93d3cubW92aXN0YXJwbHVzLmVzIiwiYXVkIjoiNDE0ZTE5MjdhMzg4NGY2OGFiYzc5ZjcyODM4MzdmZDEiLCJleHAiOjE3NDUyNzYwMzQsIm5iZiI6MTc0NDQxMjAzNH0.HdHCbmJaU-a5ASyrrygJRQF9RHj9cdN9eyx9A341H3s"}
   
   req.Header["X-Movistarplus-Deviceid"] = []string{"ea3585a776ed444d8677ad8be6ef0db3"}
   req.URL.Path = "/movistarplus/amazon.tv/sdp/mediaPlayers/ea3585a776ed444d8677ad8be6ef0db3/initData"
   
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

var body = strings.NewReader(`{"accountNumber": "00QSp000009M9gzMAC-L", "userProfile": "0", "streamMiscellanea": "HTTPS", "deviceType": "SMARTTV_OTT", "deviceManufacturerProduct": "LG", "streamDRM": "Widevine", "streamFormat": "DASH"}`)
