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
   req.URL.Host = "auth.dof6.com"
   req.URL.Path = "/movistarplus/api/devices/amazon.tv/users/authenticate"
   req.URL.Scheme = "https"
   req.Header["Authorization"] = []string{"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI3VTdlN3Y4QjhTOGg4bzlBIiwiYWNjb3VudE51bWJlciI6IjAwUVNwMDAwMDA5TTlnek1BQy1MIiwicm9sZSI6InVzZXIiLCJhcHIiOiJ3ZWJkYiIsImlzcyI6Imh0dHA6Ly93d3cubW92aXN0YXJwbHVzLmVzIiwiYXVkIjoiNDE0ZTE5MjdhMzg4NGY2OGFiYzc5ZjcyODM4MzdmZDEiLCJleHAiOjE3NDUyNzYwMzQsIm5iZiI6MTc0NDQxMjAzNH0.HdHCbmJaU-a5ASyrrygJRQF9RHj9cdN9eyx9A341H3s"}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}
