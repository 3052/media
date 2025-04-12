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
   req.URL = &url.URL{}
   req.URL.Host = "clientservices.dof6.com"
   req.URL.Path = "/movistarplus/amazon.tv/sdp/mediaPlayers/ea3585a776ed444d8677ad8be6ef0db3/initData"
   req.Method = "POST"
   req.Header["Content-Type"] = []string{"application/json"}
   value := url.Values{}
   value["qspVersion"] = []string{"ssp"}
   req.URL.RawQuery = value.Encode()
   req.Body = io.NopCloser(strings.NewReader(data))
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

const data = `
{
   "accountNumber": "00QSp000009M9gzMAC-L"
}
`
