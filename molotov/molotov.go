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
   req.Header.Set("x-molotov-agent", molotov_agent)
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
   req.Header.Set("x-molotov-agent", molotov_agent)
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

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

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
   req.Header.Set("x-molotov-agent", molotov_agent)
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
               Url string
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

func (l *Login) Playback(playUrl string) (*Playback, error) {
   req, err := http.NewRequest("", playUrl, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "138.199.15.158")
   req.Header.Set("x-molotov-agent", molotov_agent)
   req.URL.RawQuery = url.Values{
      "access_token": {l.AccessToken},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   play := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   for key, value := range p.UpDrm.License.HttpHeaders {
      req.Header.Set(key, value)
   }
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

func (p *Playback) FhdReady() string {
   return strings.Replace(p.Stream.Url, "high", "fhdready", 1)
}

type Playback struct {
   Stream struct {
      Url string // MPD
   }
   UpDrm struct {
      License struct {
         HttpHeaders map[string]string `json:"http_headers"`
      }
   } `json:"up_drm"`
}
