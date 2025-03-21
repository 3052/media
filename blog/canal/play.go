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
   req.URL.Path = "/v1/assets/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi/play"
   req.URL.Scheme = "https"
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Authorization"] = []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF1dGgiOnsicyI6IncxZjhhOGZiMC0wNWZiLTExZjAtYjVkYS1mMzJkMWNkNWRkZjciLCJ1IjoiV3ozS0JWRnAzY2xwclEzZWVNUGZZZyIsImwiOiJlbl9VUyIsImQiOiJQQyIsImRtIjoiRmlyZWZveCIsIm9tIjoiTyIsImMiOiIzR01XanAwTldZT2ZhOThVZjhhbU1oUXVSNnJ6dUxvY3FSZ0NKcEZpUjI0Iiwic3QiOiJmdWxsIiwiZyI6ImV5SmljaUk2SW0wM1kzQWlMQ0prWWlJNlptRnNjMlVzSW5CMElqcG1ZV3h6WlN3aVpHVWlPaUppY21GdVpFMWhjSEJwYm1jaUxDSjFjQ0k2SW0wM1kzQWlMQ0p2Y0NJNklqRXdNRFEySW4wIiwiZiI6NiwiYiI6Im03Y3AifSwibmJmIjoxNzQyNTIzNzA5LCJleHAiOjE3NDI1MjU1NzcsImlhdCI6MTc0MjUyMzcwOSwiYXVkIjoibTdjcCJ9.lnjPwQryinqnFccT8ryVqF6joz0c_0vguiKx-w6JoiI"}
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
   "player": {
      "capabilities": {
         "mediaTypes": [ "DASH" ],
         "drmSystems": [ "Widevine" ]
      }
   }
}
`)
