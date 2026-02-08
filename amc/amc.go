package amc

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

func GetDash(sources []Source) (*Source, bool) {
   for _, source_var := range sources {
      if source_var.Type == "application/dash+xml" {
         return &source_var, true
      }
   }
   return nil, false
}

func (s *Source) Dash() (*Dash, error) {
   resp, err := http.Get(s.Src)
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

type Dash struct {
   Body []byte
   Url  *url.URL
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

func (c *Client) Refresh() error {
   req, _ := http.NewRequest("POST", "https://gw.cds.amcn.com", nil)
   req.URL.Path = "/auth-orchestration-id/api/v1/refresh"
   req.Header.Set("authorization", "Bearer "+c.Data.RefreshToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(c)
}

func (c *Client) Login(email, password string) error {
   data, err := json.Marshal(map[string]string{
      "email":    email,
      "password": password,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com", bytes.NewReader(data),
   )
   if err != nil {
      return err
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
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(c)
}

func (c *Client) SeasonEpisodes(id int) (*Node, error) {
   var data strings.Builder
   data.WriteString("/content-compiler-cr/api/v1/content/amcn/amcplus/type")
   data.WriteString("/season-episodes/id/")
   data.WriteString(strconv.Itoa(id))
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.URL = &url.URL{
      Scheme: "https", Host: "gw.cds.amcn.com", Path: data.String(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data Node
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data, nil
}

func (c *Client) SeriesDetail(id int) (*Node, error) {
   var data strings.Builder
   data.WriteString("/content-compiler-cr/api/v1/content/amcn/amcplus/type")
   data.WriteString("/series-detail/id/")
   data.WriteString(strconv.Itoa(id))
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.URL = &url.URL{
      Scheme: "https",
      Host: "gw.cds.amcn.com",
      Path: data.String(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data Node
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data, nil
}

func (c *Client) Playback(id int) (http.Header, []Source, error) {
   var link strings.Builder
   link.WriteString("https://gw.cds.amcn.com/playback-id/api/v1/playback/")
   link.WriteString(strconv.Itoa(id))
   body, err := json.Marshal(map[string]any{
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
      "POST", link.String(), bytes.NewReader(body),
   )
   if err != nil {
      return nil, nil, err
   }
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
   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   // 1 ugly
   if len(body) == 0 {
      return nil, nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
      Error string
   }
   err = json.Unmarshal(body, &result)
   if err != nil {
      return nil, nil, err
   }
   // 2 bad
   if result.Error != "" {
      return nil, nil, errors.New(result.Error)
   }
   // 3 good
   return resp.Header, result.Data.PlaybackJsonData.Sources, nil
}

type Metadata struct {
   EpisodeNumber int
   Nid           int
   Title         string
}

func (m *Metadata) String() string {
   var data strings.Builder
   if m.EpisodeNumber >= 0 {
      data.WriteString("episode = ")
      data.WriteString(strconv.Itoa(m.EpisodeNumber))
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("title = ")
   data.WriteString(m.Title)
   data.WriteString("\nnid = ")
   data.WriteString(strconv.Itoa(m.Nid))
   return data.String()
}

func (n *Node) ExtractEpisodes() ([]*Metadata, error) {
   for _, listNode := range n.Children {
      if listNode.Type != "list" {
         continue
      }
      var extractedMetadata []*Metadata
      for _, cardNode := range listNode.Children {
         if cardNode.Type == "card" && cardNode.Properties.Metadata != nil {
            extractedMetadata = append(extractedMetadata, cardNode.Properties.Metadata)
         }
      }
      return extractedMetadata, nil
   }
   return nil, errors.New("could not find episode list in the manifest")
}

type Node struct {
   Type       string `json:"type"`
   Children   []Node `json:"children,omitempty"`
   Properties struct {
      ManifestType string `json:"manifestType,omitempty"`
      Text         *struct {
         Title struct {
            Title string `json:"title"`
         } `json:"title"`
      } `json:"text,omitempty"`
      Metadata *Metadata `json:"metadata,omitempty"`
   } `json:"properties"`
}

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

type Source struct {
   KeySystems struct {
      ComWidevineAlpha *struct {
         LicenseUrl string `json:"license_url"`
      } `json:"com.widevine.alpha"`
   } `json:"key_systems"`
   Src  string // URL to the MPD manifest
   Type string // e.g., "application/dash+xml"
}

func (s *Source) Widevine(header http.Header, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", s.KeySystems.ComWidevineAlpha.LicenseUrl,
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
