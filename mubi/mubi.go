package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

var ClientCountry = "US"

func join(items ...string) string {
   return strings.Join(items, "")
}

// if you do this you get 4K, but ONLY VP9 no H.264
var ClientAcceptVideoCodecs = "vp9"

func (s *Session) SecureUrl(filmId int) (*SecureUrl, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host: "api.mubi.com",
      Path: join("/v3/films/", strconv.Itoa(filmId), "/viewing/secure_url"),
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("user-agent", "Firefox")
   if ClientAcceptVideoCodecs != "" {
      req.Header.Set("client-accept-video-codecs", ClientAcceptVideoCodecs)
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SecureUrl
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.UserMessage != "" {
      return nil, errors.New(result.UserMessage)
   }
   return &result, nil
}

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

func (s *SecureUrl) Mpd() (*Mpd, error) {
   resp, err := http.Get(s.Url)
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

const forbidden = "HTTP Status 403 – Forbidden"

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
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

func FetchId(slug string) (int, error) {
   req, _ := http.NewRequest("", "https://api.mubi.com", nil)
   req.URL.Path = "/v3/films/" + slug
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return 0, err
   }
   defer resp.Body.Close()
   var result struct {
      Id int
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return 0, err
   }
   return result.Id, nil
}

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

func (l *LinkCode) String() string {
   var data strings.Builder
   data.WriteString("TO LOG IN AND START WATCHING\n")
   data.WriteString("Go to\n")
   data.WriteString("mubi.com/android\n")
   data.WriteString("and enter the code below\n")
   data.WriteString(l.LinkCode)
   return data.String()
}

func (l *LinkCode) Session() (*Session, error) {
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
   result := &Session{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Session struct {
   Token string
   User  struct {
      Id int
   }
}

func (s *Session) Widevine(data []byte) ([]byte, error) {
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
      "sessionId": s.Token,
      "userId":    s.User.Id,
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
   var result struct {
      License []byte
   }
   err = json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (s *Session) Viewing(filmId int) error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "api.mubi.com",
      Path: join("/v3/films/", strconv.Itoa(filmId), "/viewing"),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      UserMessage string `json:"user_message"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   if result.UserMessage != "" {
      return errors.New(result.UserMessage)
   }
   return nil
}
