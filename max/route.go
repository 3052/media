package max

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

func (r *routes) edit() (*edit, bool) {
   for _, included := range r.Included {
      if included.Relationships != nil {
         if included.Relationships.Edit != nil {
            return included.Relationships.Edit, true
         }
      }
   }
   return nil, false
}

type edit struct {
   Data struct {
      Id string
   }
}

type routes struct {
   Data struct {
      Attributes struct {
         Url string
      }
   }
   Included []struct {
      Relationships *struct {
         Show *struct {
            Data struct {
               Id string
            }
         }
         Edit *edit
      }
   }
}

func (n Login) routes(route string) (*routes, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.RawQuery = url.Values{
      "include": {"default"},
      // this is not required, but results in a smaller response
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer " + n.Data.Attributes.Token)
   req.URL.Path = "/cms/routes" + route
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   routes1 := &routes{}
   err = json.NewDecoder(resp.Body).Decode(routes1)
   if err != nil {
      return nil, err
   }
   return routes1, nil
}
