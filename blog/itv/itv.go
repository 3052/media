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
   value["variables"] = []string{variables}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

///

const query = `
query ProgrammePage(
   $brandLegacyId: BrandLegacyId
) {
  titles(
    filter: {
      brandLegacyId: $brandLegacyId
      available: "NOW"
      platform: MOBILE
      tiers: ["FREE", "PAID"]
    }
  ) {
    ... on Title {
     latestAvailableVersion {
       playlistUrl
     }
   }
  }
}
`

const variables = `
{
  "broadcaster": "ITV",
  "brandLegacyId": "18910",
  "features": [
    "HD",
    "SINGLE_TRACK",
    "MPEG_DASH",
    "WIDEVINE",
    "WIDEVINE_DOWNLOAD",
    "INBAND_TTML",
    "OUTBAND_WEBVTT",
    "INBAND_AUDIO_DESCRIPTION"
  ]
}
`
