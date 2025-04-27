package amc

import (
   "bufio"
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
)

func (a *Auth) Playback(id int) (Byte[Playback], error) {
   data, err := json.Marshal(map[string]any{
      "adtags": map[string]any{
         "lat":          0,
         "mode":         "on-demand",
         "playerHeight": 0,
         "playerWidth":  0,
         "ppid":         0,
         "url":          "-",
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/playback-id/api/v1/playback/" + strconv.Itoa(id)
   req.Header = http.Header{
      "authorization":       {"Bearer " + a.Data.AccessToken},
      "content-type":        {"application/json"},
      "x-amcn-device-ad-id": {"-"},
      "x-amcn-language":     {"en"},
      "x-amcn-network":      {"amcplus"},
      "x-amcn-platform":     {"web"},
      "x-amcn-service-id":   {"amcplus"},
      "x-amcn-tenant":       {"amcn"},
      "x-ccpa-do-not-sell":  {"doNotPassData"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var buf bytes.Buffer
   err = resp.Write(&buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

type Auth struct {
   Data struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

func (a *Auth) Unauth() error {
   req, _ := http.NewRequest("POST", "https://gw.cds.amcn.com", nil)
   req.URL.Path = "/auth-orchestration-id/api/v1/unauth"
   req.Header = http.Header{
      "x-amcn-device-id": {"-"},
      "x-amcn-language":  {"en"},
      "x-amcn-network":   {"amcplus"},
      "x-amcn-platform":  {"web"},
      "x-amcn-tenant":    {"amcn"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a *Auth) Refresh() (Byte[Auth], error) {
   req, _ := http.NewRequest("POST", "https://gw.cds.amcn.com", nil)
   req.URL.Path = "/auth-orchestration-id/api/v1/refresh"
   req.Header.Set("authorization", "Bearer "+a.Data.RefreshToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *Auth) Login(email, password string) (Byte[Auth], error) {
   data, err := json.Marshal(map[string]string{
      "email":    email,
      "password": password,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/auth-orchestration-id/api/v1/login"
   req.Header = http.Header{
      "authorization":           {"Bearer " + a.Data.AccessToken},
      "content-type":            {"application/json"},
      "x-amcn-device-ad-id":     {"-"},
      "x-amcn-device-id":        {"-"},
      "x-amcn-language":         {"en"},
      "x-amcn-network":          {"amcplus"},
      "x-amcn-platform":         {"web"},
      "x-amcn-service-group-id": {"10"},
      "x-amcn-service-id":       {"amcplus"},
      "x-amcn-tenant":           {"amcn"},
      "x-ccpa-do-not-sell":      {"doNotPassData"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *Auth) Unmarshal(data Byte[Auth]) error {
   return json.Unmarshal(data, a)
}

type Byte[T any] []byte

type Playback struct {
   Header http.Header
   Body   struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
   }
}

func (p *Playback) Dash() (*Source, bool) {
   for _, source1 := range p.Body.Data.PlaybackJsonData.Sources {
      if source1.Type == "application/dash+xml" {
         return &source1, true
      }
   }
   return nil, false
}

func (p *Playback) Unmarshal(data Byte[Playback]) error {
   resp, err := http.ReadResponse(
      bufio.NewReader(bytes.NewReader(data)), nil,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   p.Header = resp.Header
   return json.NewDecoder(resp.Body).Decode(&p.Body)
}

func (p *Playback) Widevine(source1 *Source, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", source1.KeySystems.Widevine.LicenseUrl, bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("bcov-auth", p.Header.Get("x-amcn-bc-jwt"))
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Source struct {
   KeySystems *struct {
      Widevine struct {
         LicenseUrl string `json:"license_url"`
      } `json:"com.widevine.alpha"`
   } `json:"key_systems"`
   Src  string // MPD
   Type string
}
