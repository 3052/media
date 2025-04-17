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
   
   //itv.com/watch/gone-girl/10a5503a0001B
   //value["variables"] = []string{` { "titleLegacyId": "10/5503/0001" } `}
   
   // itv.com/watch/goldeneye/18910
   
   value["query"] = []string{query}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

const query = `
query TitleCcidLegacyId( $titleLegacyId: TitleLegacyId ) {
   titles(
      filter: { legacyId: $titleLegacyId }
   ) { brandLegacyId }
}
`

