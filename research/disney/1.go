package disney

import (
   "encoding/json"
   "net/http"
   "net/url"
)

func (r refresh_token) explore(entity string) (*explore_page, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/explore/v1.12/page/entity-" + entity
   value := url.Values{}
   value["enhancedContainersLimit"] = []string{"1"}
   value["limit"] = []string{"1"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   req.Header.Set(
      "Authorization", "Bearer " + r.Extensions.Sdk.Token.AccessToken,
   )
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Page explore_page
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Page, nil
}

type explore_page struct {
   Actions []struct {
      ResourceId string
      Visuals struct {
         DisplayText string
      }
   }
}

func (e explore_page) restart() (string, bool) {
   for _, action := range e.Actions {
      if action.Visuals.DisplayText == "RESTART" {
         return action.ResourceId, true
      }
   }
   return "", false
}
