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
   req.URL.Host = "fapi.molotov.tv"
   req.URL.Scheme = "https"
   req.Header.Set(
      "x-molotov-agent", `{ "app_build": 4, "app_id": "browser_app" }`,
   )
   req.URL.Path = "/v3/auth/refresh/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMjgxODQxMDgiLCJleHBpcmVzIjoxNzQ4MjEyODc1LCJwcm9maWxlX2lkIjoiMjgxMzc5NjQiLCJ0b2tlbiI6IkJtdXJlN2RDRVR4amZuNmlZMHRWODcvZXFHaXdGVDZzMGdQczhGREhycW1ZRGIyUlRvak82K2x6b2YxeEFiQ0ZYVllhZ3JzWWNUa1I4TERHYy9wbEJRPT0iLCJ1c2VyX2lkIjoiMjgxODQxMDgiLCJ2IjoxfQ.P4f8QDTjjVd0yvtLI5mT50TByDXHJr-bxB9ETA1HhJI"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}
