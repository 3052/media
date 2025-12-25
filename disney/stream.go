package disney

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (a *account) obtain_license(data []byte) ([]byte, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.playback.edge.bamgrid.com"
   req.URL.Path = "/widevine/v1/obtain-license"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(bytes.NewReader(data))
   req.Header.Set("Authorization", a.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type stream struct {
   Sources []struct {
      Complete struct {
         Url string
      }
   }
}

func (a *account) stream(resource_id string) (*stream, error) {
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
      "Authorization", "Bearer "+a.AccessToken,
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
   var result struct {
      Errors []Error
      Stream stream
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Stream, nil
}
