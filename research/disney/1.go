package disney

import (
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (r *refresh_token) playback() (*playback, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.playback.edge.bamgrid.com"
   req.URL.Path = "/v7/playback/ctr-regular"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(playback_data))
   req.Header.Add("Content-Type", "application/json")
   req.Header.Add("X-Application-Version", "5d5917f8")
   req.Header.Add("X-Bamsdk-Client-Id", "disney-svod-3d9324fc")
   req.Header.Add("X-Bamsdk-Platform", "javascript/windows/firefox")
   req.Header.Add("X-Bamsdk-Version", "34.3")
   req.Header.Add("X-Dss-Feature-Filtering", "true")
   req.Header.Add(
      "Authorization", "Bearer " + r.Extensions.Sdk.Token.AccessToken,
   )
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription = ")
   data.WriteString(e.Description)
   return data.String()
}

type Error struct {
   Code string
   Description string
}

type playback struct {
   Errors []Error
   Stream struct {
      Sources []struct {
         Complete struct {
            Url string
         }
      }
   }
}

const playback_data = `
{
  "playback": {
    "attributes": {
      "assetInsertionStrategies": {
        "point": "SGAI",
        "range": "SGAI"
      }
    }
  },
  "playbackId": "eyJtZWRpYUlkIjoiYWE0MDFhMmItYjdmNC00YzExLWJmNjEtYTNiMDZmOWM5NzRkIiwiYXZhaWxJZCI6ImNkNDkwZmE0LTBkMWYtNDU1ZS04ZGNiLWZmZmQ1MTY2NmMyMSIsImF2YWlsVmVyc2lvbiI6Mywic291cmNlSWQiOiJjZDQ5MGZhNC0wZDFmLTQ1NWUtOGRjYi1mZmZkNTE2NjZjMjEiLCJjb250ZW50VHlwZSI6InZvZCJ9"
}
`
