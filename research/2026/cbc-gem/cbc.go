package cbc

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

func (c GemCatalog) Item() (*LineupItem, bool) {
   for _, content := range c.Content {
      for _, lineup := range content.Lineups {
         for _, item := range lineup.Items {
            if item.URL == c.SelectedUrl {
               return &item, true
            }
         }
      }
   }
   return nil, false
}

func (g *GemProfile) Unmarshal() error {
   return json.Unmarshal(g.Raw, &g.Gem)
}

type GemProfile struct {
   Gem struct {
      ClaimsToken string
   }
   Raw []byte
}

type MediaService struct {
   Message string
   URL string
}

func (t LoginToken) Profile() (*GemProfile, error) {
   req, err := http.NewRequest("GET", "https://services.radio-canada.ca", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/ott/subscription/v2/gem/Subscriber/profile"
   req.URL.RawQuery = "device=phone_android"
   req.Header.Set("Authorization", "Bearer " + t.Access_Token)
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   gem := new(GemProfile)
   if err := json.NewDecoder(res.Body).Decode(gem); err != nil {
      return nil, err
   }
   return gem, nil
}

func (g *GemCatalog) New(address string) error {
   // you can also use `phone_android`, but it returns combined number and name:
   // 3. Beauty Hath Strange Power
   req, err := http.NewRequest("GET", "https://services.radio-canada.ca", nil)
   if err != nil {
      return err
   }
   req.URL.RawQuery = "device=web"
   req.URL.Path, err = func() (string, error) {
      u, err := url.Parse(address)
      if err != nil {
         return "", err
      }
      return "/ott/catalog/v2/gem/show" + u.Path, nil
   }()
   if err != nil {
      return err
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(g)
}

func (p GemProfile) Media(i *LineupItem) (*MediaService, error) {
   req, err := http.NewRequest("GET", "https://services.radio-canada.ca", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/media/validation/v2"
   req.URL.RawQuery = url.Values{
      "appCode": {"gem"},
      "idMedia": {i.FormattedIdMedia},
      "manifestType": {manifest_type},
      "output": {"json"},
      // you need this one the first request for a video, but can omit after
      // that. we will just send it every time.
      "tech": {"hls"},
   }.Encode()
   req.Header = http.Header{
      "X-Claims-Token": {p.Gem.ClaimsToken},
      "X-Forwarded-For": {forwarded_for},
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   media := new(MediaService)
   if err := json.NewDecoder(res.Body).Decode(media); err != nil {
      return nil, err
   }
   if media.Message != "" {
      return nil, errors.New(media.Message)
   }
   media.URL = strings.Replace(media.URL, "[manifestType]", manifest_type, 1)
   return media, nil
}
const manifest_type = "desktop"

type LoginToken struct {
   Access_Token string
}

func (t *LoginToken) New(username, password string) error {
   address := func() string {
      var b strings.Builder
      b.WriteString("https://login.cbc.radio-canada.ca")
      b.WriteString("/bef1b538-1950-4283-9b27-b096cbc18070")
      b.WriteString("/B2C_1A_ExternalClient_ROPC_Auth/oauth2/v2.0/token")
      return b.String()
   }()
   res, err := http.PostForm(address, url.Values{
      "client_id": {"7f44c935-6542-4ce7-ae05-eb887809741c"},
      "grant_type": {"password"},
      "password": {password},
      "scope": {strings.Join(scope, " ")},
      "username": {username},
   })
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(t)
}

type LineupItem struct {
   URL string
   FormattedIdMedia string
}

const forwarded_for = "99.224.0.0"

var scope = []string{
   "https://rcmnb2cprod.onmicrosoft.com/84593b65-0ef6-4a72-891c-d351ddd50aab/subscriptions.write",
   "https://rcmnb2cprod.onmicrosoft.com/84593b65-0ef6-4a72-891c-d351ddd50aab/toutv-profiling",
   "openid",
}

