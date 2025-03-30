package plex

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a *Address) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "watch.plex.tv")
   a[0] = strings.TrimPrefix(data, "/watch")
   return nil
}

type Address [1]string

func (a Address) String() string {
   return a[0]
}

func (u User) Metadata(match1 *Match) (Byte[Metadata], error) {
   req, _ := http.NewRequest("", "https://vod.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/" + match1.RatingKey
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-token", u.AuthToken)
   if ForwardedFor != "" {
      req.Header.Set("x-forwarded-for", ForwardedFor)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (m *Metadata) Unmarshal(data Byte[Metadata]) error {
   var value struct {
      MediaContainer struct {
         Metadata []Metadata
      }
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   *m = value.MediaContainer.Metadata[0]
   return nil
}

type Metadata struct {
   Media []struct {
      Part     []Part
      Protocol string
   }
}

func (u *User) Unmarshal(data Byte[User]) error {
   return json.Unmarshal(data, u)
}

type Byte[T any] []byte

func NewUser() (Byte[User], error) {
   req, _ := http.NewRequest("POST", "https://plex.tv", nil)
   req.URL.Path = "/api/v2/users/anonymous"
   req.Header = http.Header{
      "accept":                   {"application/json"},
      "x-plex-product":           {"Plex Mediaverse"},
      "x-plex-client-identifier": {"!"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (u User) Widevine(part1 *Part, data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", part1.License, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.URL.Host = "vod.provider.plex.tv"
   req.URL.RawQuery = url.Values{
      "x-plex-drm":   {"widevine"},
      "x-plex-token": {u.AuthToken},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (u User) Mpd(part1 *Part) (*http.Response, error) {
   req, err := http.NewRequest("", part1.Key, nil)
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.URL.Host = "vod.provider.plex.tv"
   req.URL.RawQuery = "x-plex-token=" + u.AuthToken
   req.Header = http.Header{}
   if ForwardedFor != "" {
      req.Header.Set("x-forwarded-for", ForwardedFor)
   }
   return http.DefaultClient.Do(req)
}

func (m *Metadata) Dash() (*Part, bool) {
   for _, media := range m.Media {
      if media.Protocol == "dash" {
         return &media.Part[0], true
      }
   }
   return nil, false
}

var ForwardedFor string

func (u User) Match(web Address) (*Match, error) {
   req, _ := http.NewRequest("", "https://discover.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/matches"
   req.URL.RawQuery = url.Values{
      "url":          {web[0]},
      "x-plex-token": {u.AuthToken},
   }.Encode()
   req.Header.Set("accept", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      MediaContainer struct {
         Metadata []Match
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.MediaContainer.Metadata[0], nil
}

type Match struct {
   RatingKey string
}

type Part struct {
   Key     string
   License string
}

type User struct {
   AuthToken string
}
