package disney

import (
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

const playback_id = `
{
   "mediaId": "aa401a2b-b7f4-4c11-bf61-a3b06f9c974d"
}
`

func (r refresh_token) playback() (*playback, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.playback.edge.bamgrid.com"
   req.URL.Path = "/v7/playback/ctr-regular"
   req.URL.Scheme = "https"
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set(
      "Authorization", "Bearer " + r.Extensions.Sdk.Token.AccessToken,
   )
   req.Header.Set("X-Dss-Feature-Filtering", "true")
   req.Header.Set("X-Application-Version", "")
   req.Header.Set("X-Bamsdk-Client-Id", "")
   req.Header.Set("X-Bamsdk-Platform", "")
   req.Header.Set("X-Bamsdk-Version", "")
   data := base64.StdEncoding.EncodeToString([]byte(playback_id))
   data = fmt.Sprintf(`
   {
     "playback": {
       "attributes": {
         "assetInsertionStrategies": {
           "point": "SGAI",
           "range": "SGAI"
         }
       }
     },
     "playbackId": %q
   }
   `, data)
   req.Body = io.NopCloser(strings.NewReader(data))
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
