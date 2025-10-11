package cineMember

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (s *session) New() error {
   resp, err := http.Head("https://www.cinemember.nl/nl")
   if err != nil {
      return err
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "PHPSESSID" {
         s[0] = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}

func (s session) String() string {
   return s[0].String()
}

func (s *session) Set(data string) error {
   var err error
   s[0], err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

func (s session) login(email, password string) error {
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
   req.AddCookie(s[0])
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

type session [1]*http.Cookie

// must run session.login first
func (s session) stream() (*http.Response, error) {
   req, err := http.NewRequest(
      "", "https://www.cinemember.nl/elements/films/stream.php?id=917398", nil,
   )
   if err != nil {
      return nil, err
   }
   req.AddCookie(s[0])
   return http.DefaultClient.Do(req)
}
