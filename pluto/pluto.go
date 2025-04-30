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

var ForwardedFor string

type Clips struct {
   Sources []struct {
      File File
      Type string
   }
}

func (c *Clips) Dash() (*File, bool) {
   for _, source := range c.Sources {
      if source.Type == "DASH" {
         return &source.File, true
      }
   }
   return nil, false
}

// The Request's URL and Header fields must be initialized
func (f *File) Mpd() (*http.Response, error) {
   var req http.Request
   req.Method = "GET"
   req.URL = &f[0]
   req.Header = http.Header{}
   return http.DefaultClient.Do(&req)
}

type File [1]url.URL

type Vod struct {
   Episode string `json:"_id"`
   Id      string
   Name    string
   Seasons []struct {
      Episodes []Vod
   }
   Slug    string
}

func (v *Vod) Clips() (*Clips, error) {
   req, _ := http.NewRequest("", "https://api.pluto.tv", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/v2/episodes/")
      if v.Id != "" {
         b.WriteString(v.Id)
      } else {
         b.WriteString(v.Episode)
      }
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

type Address [2]string

func (a *Address) Set(data string) error {
   for {
      var (
         key string
         ok  bool
      )
      key, data, ok = strings.Cut(data, "/")
      if !ok {
         return nil
      }
      switch key {
      case "movies":
         a[0] = data
      case "series":
         a[0], data, ok = strings.Cut(data, "/")
         if !ok {
            return errors.New("episode")
         }
      case "episode":
         a[1] = data
      }
   }
}

func (a *Address) Vod() (*Vod, error) {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {"web"},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "drmCapabilities":   {"widevine:L3"},
      "seriesIDs":         {a[0]},
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
   vod1 := value.Vod[0]
   if vod1.Slug != a[0] {
      if vod1.Id != a[0] {
         return nil, errors.New(vod1.Slug)
      }
   }
   for _, season1 := range vod1.Seasons {
      for _, episode := range season1.Episodes {
         if episode.Episode == a[1] {
            return &episode, nil
         }
         if episode.Slug == a[1] {
            return &episode, nil
         }
      }
   }
   return &vod1, nil
}
