package molotov

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (l *Login) ProgramView(media *MediaId) (*ProgramView, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL.RawQuery = url.Values{
      "access_token": {l.Auth.AccessToken},
   }.Encode()
   req.URL.Path = fmt.Sprint(
      "/v2/channels/", media.Channel,
      "/programs/", media.Program,
      "/view",
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &ProgramView{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   if result.Program.Actions.Play == nil {
      return nil, errors.New("program is not available for playback")
   }
   return result, nil
}

type MediaId struct {
   Program int
   Channel int
}

// https://molotov.tv/fr_fr/p/15082-531
// https://molotov.tv/fr_fr/p/15082-531/la-vie-aquatique
func (m *MediaId) Parse(link string) error {
   // Cut everything before the ID section "/p/"
   _, segment, found := strings.Cut(link, "/p/")
   if !found {
      return errors.New("invalid Molotov URL: missing '/p/' marker")
   }
   // fmt.Sscanf parses the numbers into the struct fields.
   // It stops reading Channel at the first non-digit (like the '/' in a slug).
   _, err := fmt.Sscanf(segment, "%d-%d", &m.Program, &m.Channel)
   if err != nil {
      return fmt.Errorf("failed to parse program and channel IDs: %w", err)
   }
   return nil
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

func (a *Asset) Mpd() (*Mpd, error) {
   resp, err := http.Get(strings.Replace(a.Stream.Url, "high", "fhdready", 1))
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

type Asset struct {
   Drm struct {
      Token string
   }
   Error  *AssetError
   Stream struct {
      Url string // MPD
   }
}

func (a *Asset) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-dt-auth-token", a.Drm.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

func (l *Login) Asset(view *ProgramView) (*Asset, error) {
   req, err := http.NewRequest("", view.Program.Actions.Play.Url, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "138.199.15.158")
   req.Header.Set("x-molotov-agent", browser_app)
   query := req.URL.Query() // keep existing query string
   query.Set("access_token", l.Auth.AccessToken)
   req.URL.RawQuery = query.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Asset
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != nil {
      return nil, result.Error
   }
   return &result, nil
}

func (a *AssetError) Error() string {
   var data strings.Builder
   data.WriteString("developer message = ")
   data.WriteString(a.DeveloperMessage)
   data.WriteString("\nuser message = ")
   data.WriteString(a.UserMessage)
   return data.String()
}

type AssetError struct {
   DeveloperMessage string `json:"developer_message"`
   UserMessage      string `json:"user_message"`
}

type ProgramView struct {
   Program struct {
      Actions struct {
         Play *struct { // FIXME check for nil
            Url string // fapi.molotov.tv/v2/me/assets
         }
      }
   }
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (l *Login) Refresh() error {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + l.Auth.RefreshToken
   req.Header.Set("x-molotov-agent", customer_area)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

type Login struct {
   Auth struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

func (l *Login) Fetch(email, password string) error {
   data, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email":      email,
      "password":   password,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://fapi.molotov.tv/v3.1/auth/login",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("x-molotov-agent", customer_area)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

const (
   browser_app   = `{ "app_build": 4, "app_id": "browser_app", "inner_app_version_name": "5.7.0" }`
   customer_area = `{ "app_build": 1, "app_id": "customer_area" }`
)
