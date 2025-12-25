package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strings"
)

func (r refresh_token) playback(resource_id string) (*playback, error) {
   data, err := json.Marshal(map[string]any{
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategies": map[string]string{
               "point": "SGAI",
               "range": "SGAI",
            },
         },
      },
      "playbackId": resource_id,
   })
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest(
      "POST", "https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular",
      bytes.NewReader(data),
   )
   req.Header.Set(
      "Authorization", "Bearer " + r.Extensions.Sdk.Token.AccessToken,
   )
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("X-Application-Version", "")
   req.Header.Set("X-Bamsdk-Client-Id", "")
   req.Header.Set("X-Bamsdk-Platform", "")
   req.Header.Set("X-Bamsdk-Version", "")
   req.Header.Set("X-Dss-Feature-Filtering", "true")
   resp, err := http.DefaultClient.Do(req)
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
