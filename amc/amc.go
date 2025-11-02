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

// FIXME
func (n *Node) Episodes() iter.Seq[*Node] {
   return func(yield func(*Node) bool) {
      for _, child1 := range n.Children {
         for _, child2 := range child1.Children {
            if !yield(child2) {
               return
            }
         }
      }
   }
}

type Node struct {
   Type       string
   Children   []*Node
   Properties *struct {
      Text *struct {
         Title struct {
            Title string
         }
      }
      Metadata *Metadata
   }
}

func (m *Metadata) String() string {
   var b []byte
   if m.EpisodeNumber >= 0 {
      b = []byte("episodeNumber = ")
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

func (n *Node) ExtractSeasons() ([]*Metadata, error) {
   seasonsTabNode, found := n.findSeasonsTabNode()
   if !found {
      return nil, errors.New("could not find the 'Seasons' tab in the JSON data")
   }
   for _, childNode := range seasonsTabNode.Children {
      if childNode.Type == "tab_bar" {
         seasonsList := childNode.Children
         extractedMetadata := make([]*Metadata, 0, len(seasonsList))
         for _, seasonNode := range seasonsList {
            if seasonNode.Properties != nil {
               if seasonNode.Properties.Metadata != nil {
                  extractedMetadata = append(extractedMetadata, seasonNode.Properties.Metadata)
               }
            }
         }
         return extractedMetadata, nil
      }
   }
   return nil, errors.New("could not find the list of seasons inside the 'Seasons' tab")
}

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
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

type Metadata struct {
   EpisodeNumber int64
   Nid           int64
   Title         string
}

func (c *Client) SeriesDetail(id int64) (*Node, error) {
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
      Data Node
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (c *Client) SeasonEpisodes(id int64) (*Node, error) {
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
      Data Node
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (n *Node) findSeasonsTabNode() (*Node, bool) {
   for _, topLevelChild := range n.Children {
      if topLevelChild.Type == "tab_bar" {
         for _, tabItem := range topLevelChild.Children {
            if tabItem.Properties != nil {
               if tabItem.Properties.Text != nil {
                  if tabItem.Properties.Text.Title.Title == "Seasons" {
                     return tabItem, true
                  }
               }
            }
         }
      }
   }
   return nil, false
}
