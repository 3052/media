package paramount

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func login(at, username, password string) (*http.Response, error) {
   data := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "www.paramountplus.com",
      Path: "/apps-api/v2.0/androidphone/auth/login.json",
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   req.Header = http.Header{}
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("user-agent", "!")
   req.Body = io.NopCloser(strings.NewReader(data))
   return http.DefaultClient.Do(&req)
}
