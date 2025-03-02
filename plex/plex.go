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

var ForwardedFor string

func (u *User) New() error {
   req, _ := http.NewRequest("POST", "https://plex.tv", nil)
   req.URL.Path = "/api/v2/users/anonymous"
   req.Header = http.Header{
      "accept": {"application/json"},
      "x-plex-product": {"Plex Mediaverse"},
      "x-plex-client-identifier": {"!"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(u)
}

type Address [1]string

func (a Address) String() string {
   return a[0]
}

func (a *Address) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "watch.plex.tv")
   (*a)[0] = strings.TrimPrefix(data, "/watch")
   return nil
}

func (u User) Match(web Address) (*Match, error) {
   req, _ := http.NewRequest("", "https://discover.provider.plex.tv", nil)
   req.URL.Path = "/library/metadata/matches"
   req.URL.RawQuery = url.Values{
      "url": {web[0]},
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

type User struct {
   AuthToken string
}

type Match struct {
   RatingKey string
}

type Part struct {
   Key string
   License string
}

type Metadata struct {
   Media []struct {
      Part []Part
      Protocol string
   }
}

func (u User) Metadata(match1 *Match) (*Metadata, error) {
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
   var value struct {
      MediaContainer struct {
         Metadata []Metadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.MediaContainer.Metadata[0], nil
}

///

type Client struct {
   AuthToken string
   Part Part
}

func (m *Metadata) Dash(user1 User) (*Client, bool) {
   for _, media := range m.Media {
      if media.Protocol == "dash" {
         var c Client
         c.AuthToken = user1.AuthToken
         c.Part = media.Part[0]
         return &c, true
      }
   }
   return nil, false
}

func (c *Client) Mpd() (*http.Response, error) {
   req, err := http.NewRequest("", c.Part.Key, nil)
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.URL.Host = "vod.provider.plex.tv"
   req.URL.RawQuery = url.Values{
      "x-plex-token": {c.AuthToken},
   }.Encode()
   req.Header = http.Header{}
   if ForwardedFor != "" {
      req.Header.Set("x-forwarded-for", ForwardedFor)
   }
   return http.DefaultClient.Do(req)
}

func (c *Client) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", c.Part.License, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.URL.Host = "vod.provider.plex.tv"
   req.URL.RawQuery = url.Values{
      "x-plex-drm": {"widevine"},
      "x-plex-token": {c.AuthToken},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
