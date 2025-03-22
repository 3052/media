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
   "ssoToken": "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2V5IjoibTcifQ..2DuD6BzA-bRjSJL3ZkkfJg.0twuBH3-p5Tnhmlqpm5R40VivdR99_5ef75wmVVY9aMHT1Mkehc9SlAzVXZTBxJJRuzzwnIowb63b4pr-cPbmKAG7u96NYJ1aS0-pYbCpMYWWRfyRZh1vrBo71SLJ4KOchauhy3utgHEpzVZwLu66DqcapBWlZ5mEyxsFnH18-X36IMYKC0qKcUA2X0VWQPpAhYaARyLhxF4QuEbiXuLSvDz9pckRxOBQfAm3lC2fYf7bFtOV0wwL95N0kJBCjE9LKVW5gFweXOuSawmcaYU1XPMHTX-M3JrxGVfXYUBZUq6AT8otUoLh3LD8IYnRPl-bbAM7u505YIoqrVnkMaPX9XGyn9ScrY5hiOpWFFvILPGdYFbInZhqHpIsSLX5A-Qr9xDMs_VMFq-6HfSLAshww.z9rdppcbL7vxO6hfTqOJrA"
}
`)
