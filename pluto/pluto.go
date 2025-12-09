package pluto

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
   "strings"
)

func NewAlfa(id string) (*Alfa, error) {
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
      Vod []Alfa
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

type Alfa struct {
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

func (a *Alfa) String() string {
   var (
      data     []byte
      lines bool
   )
   for _, season := range a.Seasons {
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

type Bravo struct {
   Sources []struct {
      File File
      Type string
   }
}

func (c *Bravo) Dash() (*File, bool) {
   for _, source := range c.Sources {
      if source.Type == "DASH" {
         return &source.File, true
      }
   }
   return nil, false
}

func NewBravo(id string) (*Bravo, error) {
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
   var bravo2 []Bravo
   err = json.NewDecoder(resp.Body).Decode(&bravo2)
   if err != nil {
      return nil, err
   }
   return &bravo2[0], nil
}

type File [1]url.URL

// The Request's URL and Header fields must be initialized
func (f *File) Mpd() (*http.Response, error) {
   var req http.Request
   req.Method = "GET"
   req.URL = &f[0]
   req.Header = http.Header{}
   return http.DefaultClient.Do(&req)
}

// these return a valid response body, but response status is "403 OK":
// http://siloh-fs.plutotv.net
// http://siloh-ns1.plutotv.net
// https://siloh-fs.plutotv.net
// https://siloh-ns1.plutotv.net
func (f *File) UnmarshalText(data []byte) error {
   err := f[0].UnmarshalBinary(data)
   if err != nil {
      return err
   }
   f[0].Scheme = "http"
   f[0].Host = "silo-hybrik.pluto.tv.s3.amazonaws.com"
   return nil
}
