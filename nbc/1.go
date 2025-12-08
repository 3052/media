package nbc

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
)

func FetchMetadata(name int) (*Metadata, error) {
   variables, err := json.Marshal(map[string]any{
      "app":      "nbc",
      "name":     strconv.Itoa(name),
      "oneApp":   true,
      "platform": "android",
      "type":     "VIDEO",
      "userId":   "",
   })
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://friendship.nbc.com/v3/graphql", nil)
   req.URL.RawQuery = url.Values{
      "query": {graphql_compact(bonanza_page)},
      "variables": {string(variables)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var body struct {
      Data struct {
         BonanzaPage struct {
            Metadata Metadata
         }
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&body)
   if err != nil {
      return nil, err
   }
   if err := body.Errors; len(err) >= 1 {
      return nil, errors.New(err[0].Message)
   }
   return &body.Data.BonanzaPage.Metadata, nil
}
