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

type MediaPart struct {
   Key     string
   License string
}

// https://watch.plex.tv/movie/memento-2000
// https://watch.plex.tv/watch/movie/memento-2000
func GetPath(inputUrl string) (string, error) {
   u, err := url.Parse(inputUrl)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(u.Path, "/watch"), nil
}

type MediaMatch struct {
   RatingKey string
}

type ItemMetadata struct {
   Media []struct {
      Part     []MediaPart
      Protocol string
   }
}

func (a *ItemMetadata) Dash() (*MediaPart, bool) {
   for _, media := range a.Media {
      if media.Protocol == "dash" {
         return &media.Part[0], true
      }
   }
   return nil, false
}

func (u User) MediaMatch(path string) (*MediaMatch, error) {
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
   var media struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []MediaMatch
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&media)
   if err != nil {
      return nil, err
   }
   if media.Error.Message != "" {
      return nil, errors.New(media.Error.Message)
   }
   return &media.MediaContainer.Metadata[0], nil
}

func (u *User) Fetch() error {
   req, _ := http.NewRequest("POST", "https://plex.tv", nil)
   req.URL.Path = "/api/v2/users/anonymous"
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-product", "Plex Mediaverse")
   req.Header.Set("x-plex-client-identifier", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(u)
}

type User struct {
   AuthToken string
}

var ForwardedFor string

func (u User) ItemMetadata(media *MediaMatch) (*ItemMetadata, error) {
   req, _ := http.NewRequest("", "https://vod.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/" + media.RatingKey
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
   var item struct {
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&item)
   if err != nil {
      return nil, err
   }
   return &item.MediaContainer.Metadata[0], nil
}

func (u User) Widevine(part *MediaPart, data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", part.License, bytes.NewReader(data))
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

func (u User) Mpd(part *MediaPart) (*http.Response, error) {
   req, err := http.NewRequest("", part.Key, nil)
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
