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
         "app": "nbc",
         "name": "saturday-night-live/video/november-15-glen-powell/9000454161",
         "platform": "web",
         "type": "VIDEO",
         "userId": ""
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
query page(
   $app: NBCUBrands!
   $name: String!
   $platform: SupportedPlatforms!
   $type: PageType!
   $userId: String!
) {
  page(
    app: $app
    name: $name
    platform: $platform
    type: $type
    userId: $userId
  ) {
    metadata {
      ...on VideoPageMetaData {
        mpxAccountId
        mpxGuid
        programmingType
      }
    }
  }
}
`
