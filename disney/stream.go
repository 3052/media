package disney

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type stream struct {
   Sources []struct {
      Complete struct {
         Url string
      }
   }
}

func (e explore_page) restart() (string, bool) {
   for _, action := range e.Actions {
      if action.Visuals.DisplayText == "RESTART" {
         return action.ResourceId, true
      }
   }
   return "", false
}

type explore_page struct {
   Actions []struct {
      ResourceId string
      Visuals    struct {
         DisplayText string
      }
   }
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
   Code        string
   Description string
}

func (a *account) obtain_license(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/widevine/v1/obtain-license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", a.Extensions.Sdk.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

///

func (a *account) explore(entity string) (*explore_page, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.URL = &url.URL{
      Scheme: "https",
      Host: "disney.api.edge.bamgrid.com",
      Path: "/explore/v1.12/page/entity-" + entity,
      RawQuery: url.Values{
         "enhancedContainersLimit": {"1"},
         "limit": {"1"},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Page explore_page
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.Page, nil
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
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
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
