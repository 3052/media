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
   value["query"] = []string{query1}
   req.URL.Scheme = "https"
   //value["variables"] = []string{`{"hello":"10/3918/0001","broadcaster":"UNKNOWN","features":["HD","SINGLE_TRACK","MPEG_DASH","WIDEVINE","WIDEVINE_DOWNLOAD","INBAND_TTML","OUTBAND_WEBVTT","INBAND_AUDIO_DESCRIPTION"]}`}
   value["variables"] = []string{`{"hello":"18910","broadcaster":"UNKNOWN","features":["HD","SINGLE_TRACK","MPEG_DASH","WIDEVINE","WIDEVINE_DOWNLOAD","INBAND_TTML","OUTBAND_WEBVTT","INBAND_AUDIO_DESCRIPTION"]}`}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

/*
episodes need
legacyId 

movies need
brandLegacyId
*/

const query1 = `
query ProgrammePage(
  $hello: BrandLegacyId
  $broadcaster: Broadcaster
  $brandCcid: CCId
  $features: [Feature!]
) {
  titles(
    filter: {
      brandLegacyId: $hello
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
