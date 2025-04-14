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
   req.Header["Content-Length"] = []string{"64"}
   req.ContentLength = 64
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Host"] = []string{"alkasvaspub.imagenio.telefonica.net"}
   req.Header["Origin"] = []string{"https://ver.movistarplus.es"}
   req.Header["Referer"] = []string{"https://ver.movistarplus.es/"}
   req.Header["User-Agent"] = []string{"Dalvik/2.1.0 (Linux; U; Android 12; 22126RN91Y Build/SP1A.210812.016)"}
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "alkasvaspub.imagenio.telefonica.net"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(body)
   
   req.Header["X-Hzid"] = []string{"eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2N2Y1Y2NlN2FkMDg3YjI1YzBmNjRhZGIiLCJpYXQiOjE3NDQ0MTIwNDQsImlzcyI6ImVhMzU4NWE3NzZlZDQ0NGQ4Njc3YWQ4YmU2ZWYwZGIzIiwiZXhwIjoxNzQ0NDU1MjQ0fQ.cYc7fzZFKT1CU5KWxuTZtEhy6CgP0rqFDBFdyjWwyJw"}
   req.URL.Path = "/asvas/ccs/00QSp000009M9gzMAC-L/SMARTTV_OTT/ea3585a776ed444d8677ad8be6ef0db3/Session"
   
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

var body = strings.NewReader(`{"contentID":3427440,"drmMediaID":"1176568", "streamType":"AST"}`)
