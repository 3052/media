package cineMember

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func login(cookie *http.Cookie, email, password string) error {
   data := url.Values{
      "emaillogin": {email},
      "password": {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://www.cinemember.nl/elements/overlays/account/login.php",
      strings.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.AddCookie(cookie)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return err
   }
   return nil
}

// you need to bless cookie first
func stream(cookie *http.Cookie) (*http.Response, error) {
   req, err := http.NewRequest(
      "", "https://www.cinemember.nl/elements/films/stream.php?id=917398", nil,
   )
   if err != nil {
      return nil, err
   }
   req.AddCookie(cookie)
   return http.DefaultClient.Do(req)
}

func session() (*http.Cookie, error) {
   resp, err := http.Head("https://www.cinemember.nl/nl")
   if err != nil {
      return nil, err
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "PHPSESSID" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}
