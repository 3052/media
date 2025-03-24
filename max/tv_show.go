package max

import (
   "encoding/json"
   "errors"
   "iter"
   "net/http"
   "net/url"
)

func (n Login) items(route string) (items, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.RawQuery = url.Values{
      "include": {"default"},
      "page[items.size]": {"9"},
   }.Encode()
   req.Header.Set("authorization", "Bearer " + n.Data.Attributes.Token)
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
      Included items
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

func (i items) episode() iter.Seq[item] {
   return func(yield func(item) bool) {
      for _, item1 := range i {
         if item1.Attributes != nil {
            if item1.Attributes.VideoType == "EPISODE" {
               if !yield(item1) {
                  break
               }
            }
         }
      }
   }
}

type items []item

type item struct {
   Attributes *struct {
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
