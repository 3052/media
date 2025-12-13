package pluto

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

func (r *Resource) Mpd() (*url.URL, []byte, error) {
   req, err := http.NewRequest("", r.File, nil)
   if err != nil {
      return nil, nil, err
   }
   req.Host = HybrikHost
   req.URL.Host = HybrikHost
   req.URL.Scheme = HybrikScheme
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, nil, errors.New(resp.Status)
   }
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

func NewSeries(id string) (*Series, error) {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {"web"},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "drmCapabilities":   {"widevine:L3"},
      "seriesIDs":         {id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Vod []Series
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Vod[0].Id != id {
      return nil, errors.New("id mismatch")
   }
   return &result.Vod[0], nil
}

type Series struct {
   Id string
   Seasons []struct {
      Number   int64
      Episodes []struct {
         Number int64
         Name   string
         Id     string `json:"_id"`
      }
   }
}

func (s *Series) String() string {
   var (
      data     []byte
      lines bool
   )
   for _, season := range s.Seasons {
      for _, episode := range season.Episodes {
         if lines {
            data = append(data, "\n\n"...)
         } else {
            lines = true
         }
         data = append(data, "season = "...)
         data = strconv.AppendInt(data, season.Number, 10)
         data = append(data, "\nepisode = "...)
         data = strconv.AppendInt(data, episode.Number, 10)
         data = append(data, "\nname = "...)
         data = append(data, episode.Name...)
         data = append(data, "\nid = "...)
         data = append(data, episode.Id...)
      }
   }
   return string(data)
}
// these return a valid response body, but response status is "403 OK":
// http://siloh-fs.plutotv.net
// http://siloh-ns1.plutotv.net
// https://siloh-fs.plutotv.net
// https://siloh-ns1.plutotv.net
const (
   // HybrikScheme is the target protocol scheme.
   HybrikScheme = "http"
   // HybrikHost is the target host for the modified location.
   HybrikHost = "silo-hybrik.pluto.tv.s3.amazonaws.com"
)

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://service-concierge.clusters.pluto.tv/v1/wv/alt",
      "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func NewClip(id string) (*Clip, error) {
   req, _ := http.NewRequest("", "https://api.pluto.tv", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v2/episodes/")
      data.WriteString(id)
      data.WriteString("/clips.json")
      return data.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result []Clip
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result[0], nil
}

// Clip represents the top-level metadata structure.
type Clip struct {
   ID            string   `json:"_id"`
   Author        string   `json:"author"`
   Name          string   `json:"name"`
   Duration      int      `json:"duration"`
   LiveBroadcast bool     `json:"liveBroadcast"`
   Provider      string   `json:"provider"`
   Code          string   `json:"code"`
   InternalCode  string   `json:"internalCode,omitempty"`
   InPoint       int      `json:"inPoint"`
   OutPoint      int      `json:"outPoint"`
   Thumbnail     string   `json:"thumbnail"`
   Sources       []Resource 
   URL           string   `json:"url"`
   PartnerCode   string   `json:"partnerCode,omitempty"`
}

func (c *Clip) Dash() (*Resource, bool) {
   for _, source := range c.Sources {
      if source.Type == "DASH" {
         return &source, true
      }
   }
   return nil, false
}
type Resource struct {
   File       string `json:"file"`
   Type       string `json:"type"`
   Encryption string `json:"encryption"`
   ID         string `json:"_id"`
}
