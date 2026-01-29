package paramount

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func Login(at, username, password string) (*http.Cookie, error) {
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
   // randomly fails if this is missing
   req.Header.Set("user-agent", "!")
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return nil, err
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "CBS_COM" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}
