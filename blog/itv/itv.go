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
  brands(
    filter: {
      ccid: $brandCcid
      legacyId: $brandLegacyId
      tiers: ["FREE", "PAID"]
    }
  ) {
    title
    ccid
    legacyId
    tier
    imageUrl(imageType: ITVX)
    categories
    synopses {
      ninety
      epg
    }
    earliestAvailableTitle {
      ccid
    }
    latestAvailableTitle {
      ccid
    }
    latestAvailableEpisode {
      ccid
      title
    }
    series(sortBy: SEQUENCE_ASC) {
      seriesNumber
    }
    genres(filter: { hubCategory: true }) {
      name
    }
    channel {
      name
    }
    earliestAvailableSeries {
      seriesNumber
    }
    latestAvailableSeries {
      seriesNumber
      longRunning
      fullSeries
      latestAvailableEpisode {
        broadcastDateTime
        episodeNumber
      }
    }
    numberOfAvailableSeries
    visuallySigned
  }
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
fragment VariantsFields on Version {
  variants(filter: { features: $features }) {
    features
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
fragment FilmInfo on Title {
  __typename
  ... on Film {
    title
    tier
    imageUrl(imageType: ITVX)
    synopses {
      ninety
      epg
    }
    categories
    genres {
      id
      name
      hubCategory
    }
  }
}
fragment SpecialInfo on Special {
  title
  tier
  imageUrl(imageType: ITVX)
  synopses {
    ninety
    epg
  }
  categories
  genres {
    id
    name
    hubCategory
  }
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
    ...VariantsFields
    linearContent
    visuallySigned
    duration
    bsl {
      playlistUrl
    }
  }
  ...TitleAttributionFragment
  ... on Episode {
    __typename
    ...EpisodeInfo
  }
  ... on Film {
    __typename
    ...FilmInfo
  }
  ... on Special {
    __typename
    ...SpecialInfo
  }
}
`
