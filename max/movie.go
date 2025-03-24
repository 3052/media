package max

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

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

func (n Login) movie(route string) ([]movie_item, error) {
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
   var value struct {
      Errors []struct {
         Detail string
      }
      Included []movie_item
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Detail)
   }
   return value.Included, nil
}
