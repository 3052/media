package cbc

import (
   "2a.pages.dev/rosso/http"
   "encoding/json"
   "net/url"
)

func (a Asset) Name() string {
   var b []byte
   b = append(b, a.Series...)
   if a.Episode >= 1 {
      b = append(b, sep_big...)
      b = append(b, 'S')
      b = strconv.AppendInt(b, a.Season, 10)
      b = append(b, sep_small)
      b = append(b, 'E')
      b = strconv.AppendInt(b, a.Episode, 10)
      b = append(b, sep_big...)
      b = append(b, a.Title...)
   } else {
      b = append(b, sep_big...)
      b = append(b, a.Credits.Release_Date...)
   }
   return string(b)
}

const (
   sep_big = " - "
   sep_small = ' '
)

func new_catalog_gem(link string) (*catalog_gem, error) {
   // you can also use `phone_android`, but it returns combined number and name:
   // 3. Beauty Hath Strange Power
   req := http.Get(&url.URL{
      Scheme: "https",
      Host: "services.radio-canada.ca",
      Path: "/ott/catalog/v2/gem/show/" + link,
      RawQuery: "device=web",
   })
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   gem := new(catalog_gem)
   if err := json.NewDecoder(res.Body).Decode(gem); err != nil {
      return nil, err
   }
   return gem, nil
}

func (c catalog_gem) item() *lineup_item {
   for _, content := range c.Content {
      for _, lineup := range content.Lineups {
         for _, item := range lineup.Items {
            if item.URL == c.Selected_URL {
               return &item
            }
         }
      }
   }
   return nil
}

type catalog_gem struct {
   Selected_URL string `json:"selectedUrl"`
   Content []struct {
      Lineups []struct {
         Items []lineup_item
      }
   }
   Structured_Metadata metadata `json:"structuredMetadata"`
}

type lineup_item struct {
   URL string
   ID_Media int `json:"idMedia"`
}

type metadata struct {
   Part_Of_Series *struct {
      Name string
   } `json:"partofSeries"`
   Part_Of_Season *struct {
      Season_Number int `json:"seasonNumber"`
   } `json:"partofSeason"`
   Episode_Number *int `json:"episodeNumber"`
   Name string
   Date_Created string `json:"dateCreated"`
}
