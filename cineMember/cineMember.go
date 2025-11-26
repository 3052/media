package cineMember

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
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

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      if path.Ext(req.URL.Path) != ".m4s" {
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
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

func (s Session) String() string {
   return s[0].String()
}

func (s *Session) Set(data string) error {
   var err error
   s[0], err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

type Stream struct {
   Error string
   Links []struct {
      MimeType string
      Url      string
   }
   NoAccess bool
}

func (s *Stream) Vtt() (string, bool) {
   for _, link := range s.Links {
      if link.MimeType == "text/vtt" {
         return link.Url, true
      }
   }
   return "", false
}

func (s *Stream) Dash() (string, bool) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return link.Url, true
      }
   }
   return "", false
}

// must run Session.Login first
func (s Session) Stream(id int) (*Stream, error) {
   req, _ := http.NewRequest("", "https://www.cinemember.nl", nil)
   req.URL.Path = "/elements/films/stream.php"
   req.URL.RawQuery = "id=" + fmt.Sprint(id)
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var stream_var Stream
   err = json.NewDecoder(resp.Body).Decode(&stream_var)
   if err != nil {
      return nil, err
   }
   if stream_var.Error != "" {
      return nil, errors.New(stream_var.Error)
   }
   if stream_var.NoAccess {
      return nil, errors.New("no access")
   }
   return &stream_var, nil
}
