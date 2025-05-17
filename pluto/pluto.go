package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type SeriesId string

// pluto.tv/on-demand/movies/623a01faef11000014cf41f7
// pluto.tv/on-demand/movies/623a01faef11000014cf41f7/details
// pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6
// pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6/season/1
func (s *SeriesId) Set(data string) error {
   for {
      var (
         before string
         found  bool
      )
      before, data, found = strings.Cut(data, "/")
      if !found {
         return errors.New(`"/" not found`)
      }
      switch before {
      case "movies", "series":
         before, _, _ = strings.Cut(data, "/")
         *s = SeriesId(before)
         return nil
      }
   }
}

func (s SeriesId) Vod() (*Vod, error) {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {"web"},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "drmCapabilities":   {"widevine:L3"},
      "seriesIDs":         {string(s)},
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

///

func (v Vod) String() string {
   return ""
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

func (s SeriesId) Clips() (*Clips, error) {
   req, _ := http.NewRequest("", "https://api.pluto.tv", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/v2/episodes/")
      b.WriteString(string(s))
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
