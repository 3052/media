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

func (u User) Dash(part *MediaPart, forwardedFor string) (*Dash, error) {
   var req http.Request
   req.Header = http.Header{}
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "vod.provider.plex.tv",
      Path:     part.Key, // /library/parts/6730016e43b96c02321d7860-dash.mpd
      RawQuery: "x-plex-token=" + u.AuthToken,
   }
   resp, err := http.DefaultClient.Do(&req)
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

type ItemMetadata struct {
   Media []struct {
      Part     []MediaPart
      Protocol string
   }
   RatingKey string
}

type MediaPart struct {
   Key     string
   License string
}

type User struct {
   AuthToken string
}

// https://watch.plex.tv/movie/memento-2000
// https://watch.plex.tv/watch/movie/memento-2000
func GetPath(rawUrl string) (string, error) {
   u, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(u.Path, "/watch"), nil
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

func (i *ItemMetadata) Dash() (*MediaPart, bool) {
   for _, media := range i.Media {
      if media.Protocol == "dash" {
         return &media.Part[0], true
      }
   }
   return nil, false
}
func (u *User) Fetch() error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-product", "Plex Mediaverse")
   req.Header.Set("x-plex-client-identifier", "!")
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "plex.tv",
      Path:   "/api/v2/users/anonymous",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(u)
}

func (u User) RatingKey(rawUrl string) (*ItemMetadata, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("accept", "application/json")
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "discover.provider.plex.tv",
      Path:   "/library/metadata/matches",
      RawQuery: url.Values{
         "url":          {rawUrl},
         "x-plex-token": {u.AuthToken},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error.Message != "" {
      return nil, errors.New(result.Error.Message)
   }
   return &result.MediaContainer.Metadata[0], nil
}

func (u User) Media(item *ItemMetadata, forwardedFor string) (*ItemMetadata, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-token", u.AuthToken)
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/metadata/" + item.RatingKey,
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
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.MediaContainer.Metadata[0], nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
