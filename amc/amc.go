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
      return nil, nil
   },
}

func (a *Auth) Unauth() error {
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
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a *Auth) SeasonEpisodes(id int64) (*Child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/season-episodes/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer " + a.Data.AccessToken)
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

func (a *Auth) SeriesDetail(id int64) (*Child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/series-detail/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header.Set("authorization", "Bearer " + a.Data.AccessToken)
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

func (p *Playback) Dash() (*Source, bool) {
   for _, sourceVar := range p.Body.Data.PlaybackJsonData.Sources {
      if sourceVar.Type == "application/dash+xml" {
         return &sourceVar, true
      }
   }
   return nil, false
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
