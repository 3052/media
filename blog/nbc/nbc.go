package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Add("Content-Type", "application/json")
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "friendship.nbc.com"
   req.URL.Path = "/v3/graphql"
   req.URL.Scheme = "https"
   data = fmt.Sprintf(`
   {
     "query": %q,
     "variables": {
         "userId": "",
         "name": "saturday-night-live/video/november-15-glen-powell/9000454161",
         "app": "nbc",
         "platform": "web",
         "type": "VIDEO"
     }
   }
   `, data)
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}

var data = `
fragment videoPageMetaData on VideoPageMetaData {
  mpxAccountId
  mpxGuid
  programmingType
}

query page(
  $id: ID
  $name: String!
  $queryName: QueryNames
  $type: PageType!
  $subType: PageSubType
  $nationalBroadcastType: String
  $userId: String!
  $platform: SupportedPlatforms!
  $device: String
  $profile: JSON
  $timeZone: String
  $deepLinkHandle: String
  $app: NBCUBrands!
  $nbcAffiliateName: String
  $telemundoAffiliateName: String
  $language: Languages
  $playlistMachineName: String
  $mpxGuid: String
  $authorized: Boolean
  $minimumTiles: Int
  $endCardMpxGuid: String
  $endCardTagLine: String
  $seasonNumber: Int
  $creditMachineName: String
  $roleMachineName: String
  $originatingTitle: String
  $isDayZero: Boolean
) {
  page(
    id: $id
    name: $name
    type: $type
    subType: $subType
    nationalBroadcastType: $nationalBroadcastType
    userId: $userId
    queryName: $queryName
    platform: $platform
    device: $device
    profile: $profile
    timeZone: $timeZone
    deepLinkHandle: $deepLinkHandle
    app: $app
    nbcAffiliateName: $nbcAffiliateName
    telemundoAffiliateName: $telemundoAffiliateName
    language: $language
    playlistMachineName: $playlistMachineName
    mpxGuid: $mpxGuid
    authorized: $authorized
    minimumTiles: $minimumTiles
    endCardMpxGuid: $endCardMpxGuid
    endCardTagLine: $endCardTagLine
    seasonNumber: $seasonNumber
    creditMachineName: $creditMachineName
    roleMachineName: $roleMachineName
    originatingTitle: $originatingTitle
    isDayZero: $isDayZero
  ) {
    metadata {
      ...videoPageMetaData
    }
  }
}
`
