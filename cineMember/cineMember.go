package cineMember

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// must run Session.Login first
func (s Session) Stream(id int) (*Stream, error) {
   var link strings.Builder
   link.WriteString("https://www.cinemember.nl/elements/films/stream.php?id=")
   link.WriteString(strconv.Itoa(id))
   req, _ := http.NewRequest("", link.String(), nil)
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

// extracts the numeric ID and converts it to an integer
func FetchId(link string) (int, error) {
   resp, err := http.Get(link)
   if err != nil {
      return 0, err
   }
   defer resp.Body.Close()
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return 0, err
   }
   // 1. Cut text after "app.play('"
   _, after, found := strings.Cut(data.String(), "app.play('")
   if !found {
      return 0, errors.New("start marker not found")
   }
   // 2. Cut text at the next single quote to isolate the ID string
   idStr, _, found := strings.Cut(after, "'")
   if !found {
      return 0, errors.New("closing quote not found")
   }
   // 3. Convert string to integer
   return strconv.Atoi(idStr)
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

func (m *MediaLink) Mpd() (*Mpd, error) {
   resp, err := http.Get(m.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Mpd{data, resp.Request.URL}, nil
}

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

func (s *Stream) Dash() (*MediaLink, bool) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return &link, true
      }
   }
   return nil, false
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
