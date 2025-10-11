package cineMember

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func id(address string) (int, error) {
   resp, err := http.Get(address)
   if err != nil {
      return 0, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return 0, err
   }
   _, afterMarker, found := strings.Cut(string(data), "app.play('")
   if !found {
      return 0, errors.New("could not find the start marker")
   }
   var id int
   _, err = fmt.Sscan(afterMarker, &id)
   if err != nil {
      return 0, err
   }
   return id, nil
}

// must run session.login first
func (s session) stream(id int) (*stream, error) {
   req, err := http.NewRequest(
      "", "https://www.cinemember.nl/elements/films/stream.php", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = "id=" + fmt.Sprint(id)
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var streamVar stream
   err = json.NewDecoder(resp.Body).Decode(&streamVar)
   if err != nil {
      return nil, err
   }
   if streamVar.Error != "" {
      return nil, errors.New(streamVar.Error)
   }
   return &streamVar, nil
}

type stream struct {
   Error string
   Links []struct {
      Protocol string
      Url      string
   }
}

func (s stream) mpd() (string, bool) {
   for _, link := range s.Links {
      if link.Protocol == "mpd" {
         return link.Url, true
      }
   }
   return "", false
}

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
      "password":   {password},
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
