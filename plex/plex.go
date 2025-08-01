package plex

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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

func (u User) Metadata(matchVar *Match) (Byte[Metadata], error) {
   req, _ := http.NewRequest("", "https://vod.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/" + matchVar.RatingKey
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
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
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
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-product", "Plex Mediaverse")
   req.Header.Set("x-plex-client-identifier", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (u User) Widevine(partVar *Part, data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", partVar.License, bytes.NewReader(data))
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

func (u User) Mpd(partVar *Part) (*http.Response, error) {
   req, err := http.NewRequest("", partVar.Key, nil)
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

func (u User) Match(path string) (*Match, error) {
   req, _ := http.NewRequest("", "https://discover.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/matches"
   req.URL.RawQuery = url.Values{
      "url":          {path},
      "x-plex-token": {u.AuthToken},
   }.Encode()
   req.Header.Set("accept", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []Match
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if value.Error.Message != "" {
      return nil, errors.New(value.Error.Message)
   }
   return &value.MediaContainer.Metadata[0], nil
}

func Path(data string) string {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "watch.plex.tv")
   return strings.TrimPrefix(data, "/watch")
}
