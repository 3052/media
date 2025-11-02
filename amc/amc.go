package amc

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "iter"
   "log"
   "net/http"
   "net/url"
   "strconv"
)

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
}

func (c *Client) SeriesDetail(id int64) (*Child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/series-detail/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Data Child
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (c *Client) SeasonEpisodes(id int64) (*Child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/season-episodes/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Data Child
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
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

type Client struct {
   Data struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

func (c *Client) Login(email, password string) (ClientData, error) {
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
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-group-id", "10")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type ClientData []byte

func (c *Client) Unmarshal(data ClientData) error {
   return json.Unmarshal(data, c)
}

func (c *Client) Unauth() error {
   req, _ := http.NewRequest("POST", "https://gw.cds.amcn.com", nil)
   req.URL.Path = "/auth-orchestration-id/api/v1/unauth"
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-tenant", "amcn")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(c)
}

type Metadata struct {
   EpisodeNumber int64
   Nid           int64
   Title         string
}

func (m *Metadata) String() string {
   var b []byte
   if m.EpisodeNumber >= 0 {
      b = []byte("episode = ")
      b = strconv.AppendInt(b, m.EpisodeNumber, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "title = "...)
   b = append(b, m.Title...)
   b = append(b, "\nnid = "...)
   b = strconv.AppendInt(b, m.Nid, 10)
   return string(b)
}

func (c *Client) Refresh() (ClientData, error) {
   req, _ := http.NewRequest("POST", "https://gw.cds.amcn.com", nil)
   req.URL.Path = "/auth-orchestration-id/api/v1/refresh"
   req.Header.Set("authorization", "Bearer "+c.Data.RefreshToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playback) Dash() (*Source, bool) {
   for _, sourceVar := range p.Body.Data.PlaybackJsonData.Sources {
      if sourceVar.Type == "application/dash+xml" {
         return &sourceVar, true
      }
   }
   return nil, false
}

type Child struct {
   Children   []Child
   Properties struct {
      Metadata Metadata
   }
}

func (c *Child) Episodes() iter.Seq[*Child] {
   return func(yield func(*Child) bool) {
      for _, child1 := range c.Children {
         for _, child2 := range child1.Children {
            if !yield(&child2) {
               return
            }
         }
      }
   }
}

func (c *Child) Seasons() iter.Seq[*Child] {
   return func(yield func(*Child) bool) {
      for _, child1 := range c.Children { // tab_bar
         for _, child2 := range child1.Children {
            for _, child3 := range child2.Children {
               for _, child4 := range child3.Children {
                  if !yield(&child4) {
                     return
                  }
               }
            }
         }
      }
   }
}

func (p *Playback) Widevine(sourceVar *Source, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", sourceVar.KeySystems.Widevine.LicenseUrl, bytes.NewReader(data),
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

func (c *Client) Playback(id int64) (*Playback, error) {
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
   req.URL.Path = "/playback-id/api/v1/playback/" + strconv.FormatInt(id, 10)
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var play Playback
   err = json.NewDecoder(resp.Body).Decode(&play.Body)
   if err != nil {
      return nil, err
   }
   play.Header = resp.Header
   return &play, nil
}
