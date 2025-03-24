package max

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

func (n Login) show(route string) (*show_items, error) {
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
   var show show_items
   err = json.NewDecoder(resp.Body).Decode(&show)
   if err != nil {
      return nil, err
   }
   if len(show.Errors) >= 1 {
      return nil, errors.New(show.Errors[0].Detail)
   }
   return &show, nil
}

type show_items struct {
   Errors []struct {
      Detail string
   }
   Included []show_item
}

func (s *show_items) season() (*show_item, bool) {
   for _, item := range s.Included {
      if item.Attributes != nil {
         if item.Attributes.Title == "Season" {
            return &item, true
         }
      }
   }
   return nil, false
}

type show_item struct {
   Id         string
   Attributes *struct {
      Title     string
   }
}
