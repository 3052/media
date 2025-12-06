package amc

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
)

func (n *Node) ExtractSeasons() ([]*Metadata, error) {
   for _, child := range n.Children {
      // Guard: Skip any root child that is not a tab_bar.
      if child.Type != "tab_bar" {
         continue
      }

      for _, tabItem := range child.Children {
         // Guard: Skip any tab that isn't the "Seasons" tab.
         if tabItem.Type != "tab_bar_item" {
            continue
         }
         if tabItem.Properties.Text == nil {
            continue
         }
         if tabItem.Properties.Text.Title.Title != "Seasons" {
            continue
         }

         // We've found the "Seasons" tab item. Now find the list inside it.
         for _, seasonListContainer := range tabItem.Children {
            // Guard: Skip any child that is not the tab_bar list container.
            if seasonListContainer.Type != "tab_bar" {
               continue
            }

            // Success: We found the list. Extract and return.
            seasonList := seasonListContainer.Children
            extractedMetadata := make([]*Metadata, 0, len(seasonList))
            for _, seasonNode := range seasonList {
               if seasonNode.Properties.Metadata != nil {
                  extractedMetadata = append(extractedMetadata, seasonNode.Properties.Metadata)
               }
            }
            return extractedMetadata, nil
         }
      }
   }
   // If all loops complete without returning, the target was not found.
   return nil, errors.New("could not find the seasons list within the manifest")
}

func (n *Node) ExtractEpisodes() ([]*Metadata, error) {
   for _, listNode := range n.Children {
      if listNode.Type != "list" {
         continue
      }
      list := listNode.Children
      extractedMetadata := make([]*Metadata, 0, len(list))
      for _, cardNode := range list {
         if cardNode.Type == "card" && cardNode.Properties.Metadata != nil {
            extractedMetadata = append(extractedMetadata, cardNode.Properties.Metadata)
         }
      }
      return extractedMetadata, nil
   }
   return nil, errors.New("could not find episode list in the manifest")
}

func Dash(sources []Source) (*Source, bool) {
   for _, source_var := range sources {
      if source_var.Type == "application/dash+xml" {
         return &source_var, true
      }
   }
   return nil, false
}

func (c *Client) Playback(id int64) (http.Header, []Source, error) {
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
      return nil, nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, nil, err
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
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   // 1 ugly
   if len(data) == 0 {
      return nil, nil, errors.New(resp.Status)
   }
   var value struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
      Error string
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, nil, err
   }
   // 2 bad
   if value.Error != "" {
      return nil, nil, errors.New(value.Error)
   }
   // 3 good
   return resp.Header, value.Data.PlaybackJsonData.Sources, nil
}

type Metadata struct {
   EpisodeNumber int64
   Nid           int64
   Title         string
}

func (m *Metadata) String() string {
   var data []byte
   if m.EpisodeNumber >= 0 {
      data = []byte("episode = ")
      data = strconv.AppendInt(data, m.EpisodeNumber, 10)
   }
   if data != nil {
      data = append(data, '\n')
   }
   data = append(data, "title = "...)
   data = append(data, m.Title...)
   data = append(data, "\nnid = "...)
   data = strconv.AppendInt(data, m.Nid, 10)
   return string(data)
}

type Node struct {
   Type     string `json:"type"`
   Children   []*Node `json:"children,omitempty"`
   Properties struct {
      ManifestType string `json:"manifestType,omitempty"`
      Text *struct {
         Title struct {
            Title string `json:"title"`
         } `json:"title"`
      } `json:"text,omitempty"`
      Metadata *Metadata `json:"metadata,omitempty"`
   } `json:"properties"`
}

type Source struct {
   KeySystems *struct {
      ComWidevineAlpha struct {
         LicenseUrl string `json:"license_url"`
      } `json:"com.widevine.alpha"`
   } `json:"key_systems"`
   Src  string   // URL to the MPD manifest
   Type string  // e.g., "application/dash+xml"
}

func Widevine(
   header http.Header, sourceVar *Source, data []byte,
) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", sourceVar.KeySystems.ComWidevineAlpha.LicenseUrl,
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("bcov-auth", header.Get("x-amcn-bc-jwt"))
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (c *Client) SeasonEpisodes(id int64) (*Node, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      data := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      data = append(data, "type/season-episodes/id/"...)
      data = strconv.AppendInt(data, id, 10)
      return string(data)
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

func (c *Client) SeriesDetail(id int64) (*Node, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      data := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      data = append(data, "type/series-detail/id/"...)
      data = strconv.AppendInt(data, id, 10)
      return string(data)
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

type Client struct {
   Data struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
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

///

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

func (c *Client) Unmarshal(data ClientData) error {
   return json.Unmarshal(data, c)
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
