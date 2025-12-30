package disney

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "path"
   "strings"
)

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
//const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

// ZGlzbmV5JmFwcGxlJjEuMC4w
// disney&apple&1.0.0
// const client_api_key = "ZGlzbmV5JmFwcGxlJjEuMC4w.H9L7eJvc2oPYwDgmkoar6HzhBJRuUUzt_PcaC3utBI4"

// ZGlzbmV5JmFuZHJvaWQmMS4wLjA
// disney&android&1.0.0
const client_api_key = "ZGlzbmV5JmFuZHJvaWQmMS4wLjA.bkeb0m230uUhv8qrAXuNu39tbE_mD5EEhM_NAcohjyA"

func (a *Account) Playback(playbackId string) (*Playback, error) {
   data, err := json.Marshal(map[string]any{
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategies": map[string]string{
               "point": "SGAI",
               "range": "SGAI",
            },
         },
      },
      "playbackId": playbackId,
   })
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest(
      "POST", "https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular",
      bytes.NewReader(data),
   )
   req.Header.Set("x-bamsdk-platform", "android-tv")
   req.Header.Set("x-application-version", "google")
   req.Header.Set("x-bamsdk-client-id", "disney-svod-3d9324fc")
   req.Header.Set("x-bamsdk-version", "9.10.0")
   
   //req.Header.Set("x-bamsdk-platform", "")
   //req.Header.Set("x-application-version", "")
   //req.Header.Set("x-bamsdk-client-id", "")
   //req.Header.Set("x-bamsdk-version", "")
   
   req.Header.Set("x-dss-feature-filtering", "true")
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

func (a *Account) Widevine(data []byte) ([]byte, error) {
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
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Errors []Error
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, &result.Errors[0]
   }
   return data, nil
}

func (p *Playback) Hls() (*Hls, error) {
   resp, err := http.Get(p.Stream.Sources[0].Complete.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Hls{data, resp.Request.URL}, nil
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

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func GetEntity(rawLink string) (string, error) {
   // Parse the URL to safely access its components
   link, err := url.Parse(rawLink)
   if err != nil {
      return "", err
   }
   // Get the last part of the URL path
   last_segment := path.Base(link.Path)
   // The entity might be prefixed with "entity-", so we remove it
   return strings.TrimPrefix(last_segment, "entity-"), nil
}

type Explore struct {
   Data struct {
      Errors []Error // region
      Page   struct {
         Actions []struct {
            ResourceId string
            Visuals    struct {
               DisplayText string
            }
         }
      }
   }
   Errors []Error // explore-not-supported
}

func (a *Account) Explore(entity string) (*Explore, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "disney.api.edge.bamgrid.com",
      Path:   "/explore/v1.12/page/entity-" + entity,
      RawQuery: url.Values{
         "enhancedContainersLimit": {"1"},
         "limit":                   {"1"},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Explore
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   if len(result.Data.Errors) >= 1 {
      return nil, &result.Data.Errors[0]
   }
   return &result, nil
}

type Hls struct {
   Body []byte
   Url  *url.URL
}

type Playback struct {
   Errors []Error
   Stream struct {
      Sources []struct {
         Complete struct {
            Url string
         }
      }
   }
}

func (e *Explore) PlaybackId() (string, bool) {
   for _, action := range e.Data.Page.Actions {
      switch action.Visuals.DisplayText {
      case "PLAY", "RESTART":
         return action.ResourceId, true
      }
   }
   return "", false
}
