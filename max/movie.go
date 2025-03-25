package max

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

func (n Login) movie(route string) (*movie_items, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "page[items.size]": {"2"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.Path = "/cms/routes" + route
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var movie movie_items
   err = json.NewDecoder(resp.Body).Decode(&movie)
   if err != nil {
      return nil, err
   }
   if len(movie.Errors) >= 1 {
      return nil, errors.New(movie.Errors[0].Detail)
   }
   return &movie, nil
}

type movie_items struct {
   Errors []struct {
      Detail string
   }
   Included []movie_item
}

type movie_item struct {
   Id         string
   Attributes *struct {
      Title     string
      VideoType string
   }
   Relationships *struct {
      Show *struct {
         Data struct {
            Id string
         }
      }
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}
