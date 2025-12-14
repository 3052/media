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

// Fetch performs the HEAD request to cinemember.nl and populates the Session
// with the PHPSESSID cookie.
func (s *Session) Fetch() error {
   const targetURL = "https://www.cinemember.nl/nl"
   // We ignore the error here because the method and URL are hardcoded and
   // known to be valid.
   req, _ := http.NewRequest("HEAD", targetURL, nil)
   // THIS IS NEEDED OTHERWISE SUBTITLES ARE MISSING, GOD IS DEAD
   req.Header.Set("User-Agent", "Windows")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "PHPSESSID" {
         s.Cookie = cookie
         return nil
      }
   }
   return errors.New("PHPSESSID cookie not found in response")
}

type Stream struct {
   Error    string
   Links    []MediaLink
   NoAccess bool
}

func (s *Stream) Vtt() (*MediaLink, bool) {
   for _, link := range s.Links {
      if link.MimeType == "text/vtt" {
         return &link, true
      }
   }
   return nil, false
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
   _, after, found := strings.Cut(string(data), "app.play('")
   if !found {
      return 0, errors.New("could not find the start marker")
   }
   var id int
   _, err = fmt.Sscan(after, &id)
   if err != nil {
      return 0, err
   }
   return id, nil
}

func (s *Stream) Dash() (*MediaLink, bool) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return &link, true
      }
   }
   return nil, false
}

func (m *MediaLink) Mpd() (*url.URL, []byte, error) {
   resp, err := http.Get(m.Url)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

type MediaLink struct {
   MimeType string
   Url      string
}

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
   req.AddCookie(s.Cookie)
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

// Session holds the cookie data.
type Session struct {
   Cookie *http.Cookie
}

// must run Session.Login first
func (s Session) Stream(id int) (*Stream, error) {
   req, _ := http.NewRequest("", "https://www.cinemember.nl", nil)
   req.URL.Path = "/elements/films/stream.php"
   req.URL.RawQuery = "id=" + fmt.Sprint(id)
   req.AddCookie(s.Cookie)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Stream
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   if result.NoAccess {
      return nil, errors.New("no access")
   }
   return &result, nil
}
