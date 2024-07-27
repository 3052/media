package member

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strconv"
)

type DataArticle struct {
   Data []byte
   v *struct {
      Assets         []*ArticleAsset
      CanonicalTitle string `json:"canonical_title"`
      Id             int
      Metas          []struct {
         Key   string
         Value string
      }
   }
}

func (d *DataArticle) Unmarshal() error {
   d.v = pointer(d.v)
   err := json.Unmarshal(d.Data, d.v)
   if err != nil {
      return err
   }
   for _, asset := range d.v.Assets {
      asset.article = d
   }
   return nil
}

func (DataArticle) Episode() int {
   return 0
}

func (DataArticle) Season() int {
   return 0
}

func (DataArticle) Show() string {
   return ""
}

func (d DataArticle) Title() string {
   return d.v.CanonicalTitle
}

func (d DataArticle) Film() (*ArticleAsset, bool) {
   for _, asset := range d.v.Assets {
      if asset.LinkedType == "film" {
         return asset, true
      }
   }
   return nil, false
}

func (d DataArticle) Year() int {
   for _, meta := range d.v.Metas {
      if meta.Key == "year" {
         if v, err := strconv.Atoi(meta.Value); err == nil {
            return v
         }
      }
   }
   return 0
}

func (a ArticleSlug) Article() (*DataArticle, error) {
   body, err := func() ([]byte, error) {
      var s struct {
         Query     string `json:"query"`
         Variables struct {
            ArticleUrlSlug ArticleSlug `json:"articleUrlSlug"`
         } `json:"variables"`
      }
      s.Variables.ArticleUrlSlug = a
      s.Query = query_article
      return json.Marshal(s)
   }()
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://api.audienceplayer.com/graphql/2/user",
      "application/json", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var data struct {
      Data struct {
         Article json.RawMessage
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&data)
   if err != nil {
      return nil, err
   }
   return &DataArticle{Data: data.Data.Article}, nil
}
