package cineMember

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

const query_user = `
mutation UserAuthenticate($email: String, $password: String) {
   UserAuthenticate(email: $email, password: $password) {
      access_token
   }
}
`

const query_asset = `
mutation ArticleAssetPlay($article_id: Int, $asset_id: Int) {
   ArticleAssetPlay(article_id: $article_id asset_id: $asset_id) {
      entitlements {
         ... on ArticleAssetPlayEntitlement {
            key_delivery_url
            manifest
            protocol
         }
      }
   }
}
`

const query_article = `
query Article($articleUrlSlug: String) {
   Article(full_url_slug: $articleUrlSlug) {
      ... on Article {
         assets {
            ... on Asset {
               id
               linked_type
            }
         }
         id
      }
   }
}
` // do not use `query(`

type Address [1]string

func (a Address) Article() (*Article, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_article,
      "variables": map[string]string{
         "articleUrlSlug": a[0],
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://api.audienceplayer.com/graphql/2/user",
      "application/json", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Article Article
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data.Article, nil
}

func (a Address) String() string {
   return a[0]
}

func (a *Address) Set(data string) error {
   if !strings.HasPrefix(data, "https://") {
      return errors.New("must start with https://")
   }
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   data = strings.TrimPrefix(data, "cinemember.nl")
   data = strings.TrimPrefix(data, "/nl")
   (*a)[0] = strings.TrimPrefix(data, "/")
   return nil
}

func (a *Article) Film() (*Asset, bool) {
   for _, asset1 := range a.Assets {
      if asset1.LinkedType == "film" {
         return &asset1, true
      }
   }
   return nil, false
}

type Article struct {
   Assets []Asset
   Id     int
}

type Asset struct {
   Id         int
   LinkedType string `json:"linked_type"`
}

type Byte[T any] []byte

func (e *Entitlement) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      e.KeyDeliveryUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Entitlement struct {
   KeyDeliveryUrl string `json:"key_delivery_url"`
   Manifest string // MPD
   Protocol string
}

func (p *Play) Dash() (*Entitlement, bool) {
   for _, title := range p.Data.ArticleAssetPlay.Entitlements {
      if title.Protocol == "dash" {
         return &title, true
      }
   }
   return nil, false
}

func (p *Play) Unmarshal(data Byte[Play]) error {
   err := json.Unmarshal(data, p)
   if err != nil {
      return err
   }
   if len(p.Errors) >= 1 {
      return errors.New(p.Errors[0].Message)
   }
   return nil
}

type Play struct {
   Data struct {
      ArticleAssetPlay struct {
         Entitlements []Entitlement
      }
   }
   Errors []struct {
      Message string
   }
}

func NewUser(email, password string) (Byte[User], error) {
   data, err := json.Marshal(map[string]any{
      "query": query_user,
      "variables": map[string]string{
         "email": email,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://api.audienceplayer.com/graphql/2/user",
      "application/json", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (u *User) Unmarshal(data Byte[User]) error {
   var value struct {
      Data struct {
         UserAuthenticate User
      }
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   *u = value.Data.UserAuthenticate
   return nil
}

type User struct {
   AccessToken string `json:"access_token"`
}

func (u User) Play(article1 *Article, asset1 *Asset) (Byte[Play], error) {
   data, err := json.Marshal(map[string]any{
      "query": query_asset,
      "variables": map[string]int{
         "article_id": article1.Id,
         "asset_id": asset1.Id,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.audienceplayer.com/graphql/2/user",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // need .Set to match .Get
   req.Header.Set("authorization", "Bearer " + u.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
