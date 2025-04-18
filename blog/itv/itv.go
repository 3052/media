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
   value["variables"] = []string{variables}
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

/*
itv.com/watch/joan/10a3918
itv.com/watch/goldeneye/18910
itv.com/watch/gone-girl/10a5503a0001B
*/
const variables = `
{
  "broadcaster": "UNKNOWN",
  "brandLegacyId": "10/5503/0001B",
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

const query = `
query ProgrammePage(
  $broadcaster: Broadcaster
  $brandCcid: CCId
  $brandLegacyId: BrandLegacyId
  $features: [Feature!]
) {
  titles(
    filter: {
      brandLegacyId: $brandLegacyId
      brandCCId: $brandCcid
      broadcaster: $broadcaster
      available: "NOW"
      platform: MOBILE
      features: $features
      tiers: ["FREE", "PAID"]
    }
    sortBy: SEQUENCE_ASC
  ) {
    __typename
    ...TitleFields
  }
}

fragment TitleAttributionFragment on Title {
  attribution {
    partnership {
      name
      imageUrls {
        appsRoku
      }
    }
    contentOwner {
      name
      imageUrls {
        appsRoku
      }
    }
  }
}
fragment SeriesInfo on Series {
  longRunning
  fullSeries
  seriesNumber
  numberOfAvailableEpisodes
}

fragment TitleFields on Title {
  __typename
  titleType
  ccid
  legacyId
  brandLegacyId
  title
  brand {
    ccid
    numberOfAvailableSeries
  }
  nextAvailableTitle {
    latestAvailableVersion {
      ccid
      legacyId
    }
  }
  channel {
    name
    strapline
  }
  broadcastDateTime
  synopses {
    ninety
    epg
  }
  imageUrl(imageType: ITVX)
  regionalisation
  latestAvailableVersion {
    __typename
    ccid
    legacyId
    duration
    playlistUrl
    duration
    compliance {
      displayableGuidance
    }
    availability {
      downloadable
      end
      start
      maxResolution
      adRule
    }
    linearContent
    visuallySigned
    duration
    bsl {
      playlistUrl
    }
  }
  ...TitleAttributionFragment
}
`
