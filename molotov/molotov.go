package molotov

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Playback struct {
   Error struct {
      DeveloperMessage string `json:"developer_message"`
   }
   Stream struct {
      Url string // MPD
   }
   UpDrm struct {
      KeySystems struct {
         Widevine struct {
            License struct {
               HttpHeaders map[string]string `json:"http_headers"`
            }
         }
      } `json:"key_systems"`
   } `json:"up_drm"`
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   for key, value := range p.UpDrm.KeySystems.Widevine.License.HttpHeaders {
      req.Header.Set(key, value)
   }
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.License, nil
}

const (
   browser_app = `{ "app_build": 4, "app_id": "browser_app", "inner_app_version_name": "5.7.0" }`
   customer_area = `{ "app_build": 1, "app_id": "customer_area" }`
)

func (l *Login) Playback(playUrl string) (*Playback, error) {
   req, err := http.NewRequest("", playUrl, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "138.199.15.158")
   req.Header.Set("x-molotov-agent", browser_app)
   // keep existing query string
   query := req.URL.Query()
   query.Set("access_token", l.AccessToken)
   req.URL.RawQuery = query.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var play Playback
   err = json.NewDecoder(resp.Body).Decode(&play)
   if err != nil {
      return nil, err
   }
   if play.Error.DeveloperMessage != "" {
      return nil, errors.New(play.Error.DeveloperMessage)
   }
   return &play, nil
}

func FetchLogin(email, password string) (*Login, error) {
   value := map[string]string{
      "grant_type": "password",
      "email": email,
      "password": password,
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://fapi.molotov.tv/v3.1/auth/login",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", customer_area)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value1 struct {
      Auth Login
   }
   err = json.NewDecoder(resp.Body).Decode(&value1)
   if err != nil {
      return nil, err
   }
   return &value1.Auth, nil
}

func (l *Login) Unmarshal(data LoginData) error {
   return json.Unmarshal(data, l)
}

type LoginData []byte

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (l *Login) Refresh() (LoginData, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + l.RefreshToken
   req.Header.Set("x-molotov-agent", customer_area)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type MediaId struct {
   Channel int64
   Program int64
}

// https://molotov.tv/fr_fr/p/15082-531
// https://molotov.tv/fr_fr/p/15082-531/la-vie-aquatique
func (m *MediaId) Parse(rawUrl string) error {
   _, after, found := strings.Cut(rawUrl, "/p/")
   if !found {
      return errors.New("URL does not contain the '/p/' segment")
   }
   id, _, _ := strings.Cut(after, "/")
   program, channel, found := strings.Cut(id, "-")
   if !found {
      return errors.New("ID segment: missing '-' separator")
   }
   var err error
   m.Program, err = strconv.ParseInt(program, 10, 64)
   if err != nil {
      return err
   }
   m.Channel, err = strconv.ParseInt(channel, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

type Login struct {
   AccessToken string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

func (l *Login) PlayUrl(media *MediaId) (string, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = func() string {
      data := []byte("/v2/channels/")
      data = strconv.AppendInt(data, media.Channel, 10)
      data = append(data, "/programs/"...)
      data = strconv.AppendInt(data, media.Program, 10)
      data = append(data, "/view"...)
      return string(data)
   }()
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL.RawQuery = url.Values{
      "access_token": {l.AccessToken},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var value struct {
      Program struct {
         Actions struct {
            Play *struct {
               Url string // fapi.molotov.tv/v2/me/assets
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return "", err
   }
   if value.Program.Actions.Play == nil {
      return "", errors.New("Program.Actions.Play")
   }
   return value.Program.Actions.Play.Url, nil
}

func (p *Playback) FhdReady() string {
   return strings.Replace(p.Stream.Url, "high", "fhdready", 1)
}

