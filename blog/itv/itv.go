package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header["Accept"] = []string{"multipart/mixed; deferSpec=20220824, application/json"}
   req.Header["Content-Length"] = []string{"0"}
   req.Header["User-Agent"] = []string{"ITV_Player_(Android)"}
   req.Header["X-Apollo-Operation-Id"] = []string{"f8e83859439b0a6e50ae5d6c3a1c41c39219359266afeed4f51f77d0c9588460"}
   req.Header["X-Apollo-Operation-Name"] = []string{"ProgrammePage"}
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "content-inventory.prd.oasvc.itv.com"
   req.URL.Path = "/discovery"
   value := url.Values{}
   value["operationName"] = []string{"ProgrammePage"}
   value["query"] = []string{query}
   req.URL.Scheme = "https"
   value["variables"] = []string{` { "brandLegacyId": "18910" } `}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

/*
itv.com/watch/joan/10a3918
itv.com/watch/goldeneye/18910
itv.com/watch/gone-girl/10a5503a0001B
*/
const query = `
query ProgrammePage( $brandLegacyId: BrandLegacyId ) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
   ) {
      ... on Title {
         latestAvailableVersion {
            playlistUrl
         }
      }
   }
}
`
