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

func (s *Session) Fetch() error {
   req, _ := http.NewRequest("HEAD", "https://www.cinemember.nl/nl", nil)
   // THIS IS NEEDED OTHERWISE SUBTITLES ARE MISSING, GOD IS DEAD
   req.Header.Add("user-agent", "Windows")
   resp, err := http.DefaultClient.Do(req)
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

func Id(address string) (int, error) {
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

type Session [1]*http.Cookie

func (s Session) Login(email, password string) error {
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

type StreamInfo struct {
   Error string
   Links []struct {
      MimeType string
      Url      string
   }
   NoAccess bool
}

func (s *StreamInfo) Vtt() (string, bool) {
   for _, link := range s.Links {
      if link.MimeType == "text/vtt" {
         return link.Url, true
      }
   }
   return "", false
}

func (s *StreamInfo) Dash() (string, bool) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return link.Url, true
      }
   }
   return "", false
}

// must run Session.Login first
func (s Session) StreamInfo(id int) (*StreamInfo, error) {
   req, _ := http.NewRequest("", "https://www.cinemember.nl", nil)
   req.URL.Path = "/elements/films/stream.php"
   req.URL.RawQuery = "id=" + fmt.Sprint(id)
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var info StreamInfo
   err = json.NewDecoder(resp.Body).Decode(&info)
   if err != nil {
      return nil, err
   }
   if info.Error != "" {
      return nil, errors.New(info.Error)
   }
   if info.NoAccess {
      return nil, errors.New("no access")
   }
   return &info, nil
}
