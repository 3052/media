package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "path"
   "strconv"
   "strings"
)

func (a *Authenticate) Widevine(data []byte) ([]byte, error) {
   // final slash is needed
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   data, err = json.Marshal(map[string]any{
      "merchant":  "mubi",
      "sessionId": a.Token,
      "userId":    a.User.Id,
   })
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", base64.StdEncoding.EncodeToString(data))
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if strings.Contains(string(data), forbidden[0]) {
      return nil, forbidden
   }
   var value struct {
      License []byte
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, err
   }
   return value.License, nil
}

func (s status) Error() string {
   return strings.ToLower(s[0])
}

var forbidden = status{"HTTP Status 403 – Forbidden"}

type status [1]string

type Film struct {
   Id    int64
   Title string
   Year  int
}

var ClientCountry = "US"

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (a *Authenticate) Viewing(filmVar *Film) error {
   req, _ := http.NewRequest("POST", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v3/films/")
      b = strconv.AppendInt(b, filmVar.Id, 10)
      b = append(b, "/viewing"...)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer "+a.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var value struct {
      UserMessage string `json:"user_message"`
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return err
   }
   if value.UserMessage != "" {
      return errors.New(value.UserMessage)
   }
   return nil
}

func (c *LinkCode) String() string {
   var b strings.Builder
   b.WriteString("TO LOG IN AND START WATCHING\n")
   b.WriteString("Go to\n")
   b.WriteString("mubi.com/android\n")
   b.WriteString("and enter the code below\n")
   b.WriteString(c.LinkCode)
   return b.String()
}

type Byte[T any] []byte

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}

func (a *Authenticate) Unmarshal(data Byte[Authenticate]) error {
   return json.Unmarshal(data, a)
}

type Authenticate struct {
   Token string
   User  struct {
      Id int
   }
}

func (c *LinkCode) Authenticate() (Byte[Authenticate], error) {
   data, err := json.Marshal(map[string]string{"auth_token": c.AuthToken})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.mubi.com/v3/authenticate", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func NewLinkCode() (Byte[LinkCode], error) {
   req, _ := http.NewRequest("", "https://api.mubi.com/v3/link_code", nil)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (c *LinkCode) Unmarshal(data Byte[LinkCode]) error {
   return json.Unmarshal(data, c)
}

type SecureUrl struct {
   TextTrackUrls []Text `json:"text_track_urls"`
   Url           string      // MPD
   UserMessage   string      `json:"user_message"`
}

func (s *SecureUrl) Unmarshal(data Byte[SecureUrl]) error {
   err := json.Unmarshal(data, s)
   if err != nil {
      return err
   }
   if s.UserMessage != "" {
      return errors.New(s.UserMessage)
   }
   return nil
}
func (a *Authenticate) SecureUrl(filmVar *Film) (Byte[SecureUrl], error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v3/films/")
      b = strconv.AppendInt(b, filmVar.Id, 10)
      b = append(b, "/viewing/secure_url"...)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer "+a.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (t *Text) Base() string {
   return path.Base(t.Url)
}

type Text struct {
   Id  string
   Url string
}

func (s *Slug) Parse(data string) error {
   var found bool
   _, data, found = strings.Cut(data, "/films/")
   if !found {
      return errors.New("/films/ not found")
   }
   *s = Slug(data)
   return nil
}

type Slug string

func (s Slug) Film() (*Film, error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = "/v3/films/" + string(s)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   filmVar := &Film{}
   err = json.NewDecoder(resp.Body).Decode(filmVar)
   if err != nil {
      return nil, err
   }
   return filmVar, nil
}
