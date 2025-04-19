package itv

import (
   "encoding/json"
   "net/http"
   "net/url"
)

func discovery(brand_legacy_id string) (*http.Response, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": brand_legacy_id,
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

const query = `
query ProgrammePage( $brandLegacyId: BrandLegacyId ) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
   ) {
      ... on Episode {
         series {
            seriesNumber
         }
         episodeNumber
      }
      title
      latestAvailableVersion {
         playlistUrl
      }
   }
}
`
