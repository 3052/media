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

func (u User) Mpd(part *MediaPart, forwardedFor string) (*url.URL, []byte, error) {
   req, _ := http.NewRequest("", "https://vod.provider.plex.tv", nil)
   // /library/parts/6730016e43b96c02321d7860-dash.mpd
   req.URL.Path = part.Key
   req.URL.RawQuery = "x-plex-token=" + u.AuthToken
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
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

func (u User) Media(item *ItemMetadata, forwardedFor string) (*ItemMetadata, error) {
   req, _ := http.NewRequest("", "https://vod.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/" + item.RatingKey
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-token", u.AuthToken)
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var payload struct {
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&payload)
   if err != nil {
      return nil, err
   }
   return &payload.MediaContainer.Metadata[0], nil
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

func (u User) RatingKey(rawUrl string) (*ItemMetadata, error) {
   req, _ := http.NewRequest("", "https://discover.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/matches"
   req.URL.RawQuery = url.Values{
      "url":          {rawUrl},
      "x-plex-token": {u.AuthToken},
   }.Encode()
   req.Header.Set("accept", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var item struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&item)
   if err != nil {
      return nil, err
   }
   if item.Error.Message != "" {
      return nil, errors.New(item.Error.Message)
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

func (i *ItemMetadata) Dash() (*MediaPart, bool) {
   for _, media := range i.Media {
      if media.Protocol == "dash" {
         return &media.Part[0], true
      }
   }
   return nil, false
}
