package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return nil, nil
   },
}

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

var ClientCountry = "US"

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

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

func (t *Text) Base() string {
   return path.Base(t.Url)
}

type Text struct {
   Id  string
   Url string
}

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (a *Authenticate) Viewing(film_id int64) error {
   req, _ := http.NewRequest("POST", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v3/films/")
      b = strconv.AppendInt(b, film_id, 10)
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

func (a *Authenticate) SecureUrl(film_id int64) (Byte[SecureUrl], error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v3/films/")
      b = strconv.AppendInt(b, film_id, 10)
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

// https://mubi.com/en/films/perfect-days
// https://mubi.com/en/us/films/perfect-days
// https://mubi.com/films/perfect-days
// https://mubi.com/us/films/perfect-days
func FilmSlug(address string) (string, error) {
   _, slug, found := strings.Cut(address, "/films/")
   if !found {
      return "", errors.New(`"/films/" not found in URL`)
   }
   return slug, nil
}

func FilmId(slug string) (int64, error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = "/v3/films/" + slug
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return 0, err
   }
   defer resp.Body.Close()
   var film struct {
      Id int64
   }
   err = json.NewDecoder(resp.Body).Decode(&film)
   if err != nil {
      return 0, err
   }
   return film.Id, nil
}
