package rakuten

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (a *address) episodes(season_id string) ([]episode, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {"web"},
      "market_code":       {a.market_code},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Episodes []episode
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Episodes, nil
}

func (e *episode) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(e.Title)
   b.WriteString("\nid = ")
   b.WriteString(e.Id)
   return b.String()
}

type episode struct {
   Id    string
   Title string
}
