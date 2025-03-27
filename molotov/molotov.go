package molotov

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a *assets) widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   for key, value := range a.UpDrm.License.HttpHeaders {
      req.Header.Set(key, value)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.License, nil
}

type assets struct {
   Stream struct {
      Url string // MPD
   }
   UpDrm struct {
      License struct {
         HttpHeaders map[string]string `json:"http_headers"`
      }
   } `json:"up_drm"`
}


func (a *assets) fhd_ready() string {
   return strings.Replace(a.Stream.Url, "high", "fhdready", 1)
}

func (r *refresh) assets(view1 *view) (*assets, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv/v2/me/assets", nil)
   req.URL.RawQuery = url.Values{
      "access_token": {r.AccessToken},
      "id": {view1.Program.Video.Id},
      "type": {"vod"},
   }.Encode()
   req.Header = http.Header{
      "x-forwarded-for": {"138.199.15.158"},
      "x-molotov-agent": {molotov_agent},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   assets1 := &assets{}
   err = json.NewDecoder(resp.Body).Decode(assets1)
   if err != nil {
      return nil, err
   }
   return assets1, nil
}
type view struct {
   Program struct {
      Video struct {
         Id string
      }
   }
}

func (r *refresh) view() (*view, error) {
   req, err := http.NewRequest(
      "", "https://fapi.molotov.tv/v2/channels/531/programs/15082/view", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", molotov_agent)
   req.URL.RawQuery = "access_token=" + r.AccessToken
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   view1 := &view{}
   err = json.NewDecoder(resp.Body).Decode(view1)
   if err != nil {
      return nil, err
   }
   return view1, nil
}

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

func (n *login) New(email, password string) error {
   data, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email": email,
      "password": password,
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
   req.Header.Set("x-molotov-agent", molotov_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(n)
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (r *refresh) refresh() (Byte[refresh], error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + r.RefreshToken
   req.Header.Set("x-molotov-agent", molotov_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

type refresh struct {
   AccessToken string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

type login struct {
   Auth refresh
}

func (r *refresh) unmarshal(data Byte[refresh]) error {
   return json.Unmarshal(data, r)
}
