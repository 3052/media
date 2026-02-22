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

func (l *Login) ProgramView(media *MediaId) (*ProgramView, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "fapi.molotov.tv",
      Path: join(
         "/v2/channels/",
         strconv.Itoa(media.Channel),
         "/programs/",
         strconv.Itoa(media.Program),
         "/view",
      ),
      RawQuery: url.Values{"access_token": {l.Auth.AccessToken}}.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
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

type MediaId struct {
   Program int
   Channel int
}

// Parse extracts the IDs from a Molotov URL and populates the MediaId struct.
// It maps the first number in the URL to Program and the second to Channel.
func (m *MediaId) Parse(link string) error {
   // 1. Cut the string to get everything after "/p/"
   _, remainder, found := strings.Cut(link, "/p/")
   if !found {
      return errors.New("url does not contain the /p/ marker")
   }
   // 2. Cut the remainder at the next slash to remove any trailing text (e.g., title)
   // If no slash exists, strings.Cut returns the whole string as 'idsStr', which is what we want.
   idsStr, _, _ := strings.Cut(remainder, "/")
   // 3. Cut the string at the hyphen to separate the two numbers
   progStr, chanStr, found := strings.Cut(idsStr, "-")
   if !found {
      return errors.New("invalid format: hyphen not found between IDs")
   }
   // 4. Convert and assign to the struct fields
   var err error
   m.Program, err = strconv.Atoi(progStr)
   if err != nil {
      return errors.New("program ID is not a valid integer")
   }
   m.Channel, err = strconv.Atoi(chanStr)
   if err != nil {
      return errors.New("channel ID is not a valid integer")
   }
   return nil
}
type Dash struct {
   Body []byte
   Url  *url.URL
}

func (a *Asset) Dash() (*Dash, error) {
   resp, err := http.Get(strings.Replace(a.Stream.Url, "high", "fhdready", 1))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (l *Login) Refresh() error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "fapi.molotov.tv",
      Path:   "/v3/auth/refresh/" + l.Auth.RefreshToken,
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

func join(items ...string) string {
   return strings.Join(items, "")
}
