package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
   "strings"
)

func (l *LinkCode) Fetch() error {
   req, _ := http.NewRequest("", "https://api.mubi.com/v3/link_code", nil)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

type Authenticate struct {
   Token string
   User  struct {
      Id int
   }
}

func (l *LinkCode) Authenticate() (*Authenticate, error) {
   data, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
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
   value := &Authenticate{}
   err = json.NewDecoder(resp.Body).Decode(value)
   if err != nil {
      return nil, err
   }
   return value, nil
}

func (a *Authenticate) SecureUrl(filmId int64) (*SecureUrl, error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      data := []byte("/v3/films/")
      data = strconv.AppendInt(data, filmId, 10)
      data = append(data, "/viewing/secure_url"...)
      return string(data)
   }()
   req.Header.Set("authorization", "Bearer "+a.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var secure SecureUrl
   err = json.NewDecoder(resp.Body).Decode(&secure)
   if err != nil {
      return nil, err
   }
   if secure.UserMessage != "" {
      return nil, errors.New(secure.UserMessage)
   }
   return &secure, nil
}
var ClientCountry = "US"

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

func (l *LinkCode) String() string {
   var b strings.Builder
   b.WriteString("TO LOG IN AND START WATCHING\n")
   b.WriteString("Go to\n")
   b.WriteString("mubi.com/android\n")
   b.WriteString("and enter the code below\n")
   b.WriteString(l.LinkCode)
   return b.String()
}

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url           string      // MPD
   UserMessage   string      `json:"user_message"`
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
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if strings.Contains(string(data), forbidden) {
      return nil, errors.New(strings.ToLower(forbidden))
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

const forbidden = "HTTP Status 403 – Forbidden"

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (a *Authenticate) Viewing(filmId int64) error {
   req, _ := http.NewRequest("POST", "https://api.mubi.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v3/films/")
      b = strconv.AppendInt(b, filmId, 10)
      b = append(b, "/viewing"...)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer "+a.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
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

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}
