package pluto

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func NewClips(id string) (*Clips, error) {
   req, _ := http.NewRequest("", "https://api.pluto.tv", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/v2/episodes/")
      b.WriteString(id)
      b.WriteString("/clips.json")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var clips1 []Clips
   err = json.NewDecoder(resp.Body).Decode(&clips1)
   if err != nil {
      return nil, err
   }
   return &clips1[0], nil
}

func NewVod(id string) (*Vod, error) {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {"web"},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "drmCapabilities":   {"widevine:L3"},
      "seriesIDs":         {id},
   }.Encode()
   if ForwardedFor != "" {
      req.Header.Set("x-forwarded-for", ForwardedFor)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Vod []Vod
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Vod[0], nil
}

func (v *Vod) String() string {
   var (
      b []byte
      lines bool
   )
   for _, season := range v.Seasons {
      for _, episode := range season.Episodes {
         if lines {
            b = append(b, "\n\n"...)
         } else {
            lines = true
         }
         b = append(b, "season = "...)
         b = strconv.AppendInt(b, season.Number, 10)
         b = append(b, "\nepisode = "...)
         b = strconv.AppendInt(b, episode.Number, 10)
         b = append(b, "\nname = "...)
         b = append(b, episode.Name...)
         b = append(b, "\nid = "...)
         b = append(b, episode.Id...)
      }
   }
   return string(b)
}

type Vod struct {
   Seasons []struct {
      Number   int64
      Episodes []struct {
         Number int64
         Name   string
         Id     string `json:"_id"`
      }
   }
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

var ForwardedFor string

type Clips struct {
   Sources []struct {
      File File
      Type string
   }
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

func (c *Clips) Dash() (*File, bool) {
   for _, source := range c.Sources {
      if source.Type == "DASH" {
         return &source.File, true
      }
   }
   return nil, false
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
