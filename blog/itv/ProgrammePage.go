package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "content-inventory.prd.oasvc.itv.com"
   req.URL.Path = "/discovery"
   value := url.Values{}
   value["query"] = []string{query}
   
   // itv.com/watch/goldeneye/18910
   //value["variables"] = []string{` { "brandLegacyId": "18910" } `}
   
   //itv.com/watch/gone-girl/10a5503a0001B
   //value["variables"] = []string{` { "brandLegacyId": "10_5503_0001B" } `}
   value["variables"] = []string{` { "brandLegacyId": "10/5503/0001B" } `}
   
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

const query1 = `
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
