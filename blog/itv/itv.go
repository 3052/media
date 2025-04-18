package itv

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

const query = `
query ProgrammePage(
   $brandLegacyId: BrandLegacyId
) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
   ) {
      ... on Title {
         latestAvailableVersion { playlistUrl }
      }
   }
}
`

type legacy_id [1]string

func (v legacy_id) String() string {
   return v[0]
}

// itv.com/watch/gone-girl/10a5503a0001B
// itv.com/watch/gone-girl/10_5503_0001B
func (v *legacy_id) Set(data string) error {
   data = strings.ReplaceAll(data, "_", "/")
   v[0] = strings.ReplaceAll(data, "a", "/")
   return nil
}

func (v legacy_id) programme_page() (*http.Response, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": v[0],
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "", "https://content-inventory.prd.oasvc.itv.com/discovery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "query": {query},
      "variables": {string(data)},
   }.Encode()
   return http.DefaultClient.Do(req)
}
