package disney

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func (e explore_page) restart() (string, bool) {
   for _, action := range e.Actions {
      if action.Visuals.DisplayText == "RESTART" {
         return action.ResourceId, true
      }
   }
   return "", false
}

type explore_page struct {
   Actions []struct {
      ResourceId string
      Visuals    struct {
         DisplayText string
      }
   }
}

func (a *account) explore(entity string) (*explore_page, error) {
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
      "Authorization", "Bearer " + a.AccessToken,
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
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.Page, nil
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription = ")
   data.WriteString(e.Description)
   return data.String()
}

type Error struct {
   Code        string
   Description string
}
