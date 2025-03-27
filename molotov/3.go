package molotov

import (
   "bytes"
   "strings"
   "encoding/json"
   "net/http"
   "net/url"
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
