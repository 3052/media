package nbc

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
)

// saturday-night-live/video/november-15-glen-powell/9000454161
func FetchMetadata(name string) (*Metadata, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_page,
      "variables": map[string]string{
         "app": "nbc",
         "name": name,
         "platform": "web",
         "type": "VIDEO",
         "userId": "",
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://friendship.nbc.com/v3/graphql", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var body struct {
      Data struct {
         Page struct {
            Metadata Metadata
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&body)
   if err != nil {
      return nil, err
   }
   return &body.Data.Page.Metadata, nil
}

const query_page = `
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
