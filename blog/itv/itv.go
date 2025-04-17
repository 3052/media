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
fragment EpisodeInfo on Episode {
  series {
    __typename
    ...SeriesInfo
  }
  episodeNumber
  tier
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
}
`
