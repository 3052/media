package draken

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func FetchAsset(customId string) (*AssetItem, error) {
   data, err := json.Marshal(map[string]any{
      "query": get_custom_id,
      "variables": map[string]string{
         "customId": customId,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/apiql/v2",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   // this value is important, with the wrong value you get random failures
   req.Header.Set("x-forwarded-for", "95.192.0.0")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Viewer struct {
            ViewableCustomId *AssetItem
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Data.Viewer.ViewableCustomId == nil {
      return nil, errors.New("ViewableCustomId")
   }
   return result.Data.Viewer.ViewableCustomId, nil
}

type Entitlement struct {
   Error *struct {
      UserMessage string `json:"user_message"`
   }
   Token string
}

type Playback struct {
   Headers  map[string]string
   Playlist string // MPD
}

type AssetItem struct {
   DefaultPlayable struct {
      Id string
   }
}

const get_custom_id = `
query GetCustomIdFullMovie($customId: ID!) {
   viewer {
      viewableCustomId(customId: $customId) {
         ... on Movie {
            defaultPlayable {
               id
            }
         }
      }
   }
}
`

type Login struct {
   Token string
}

func (l Login) Widevine(play *Playback, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/playback/v1/widevine/license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   for key, value := range play.Headers {
      req.Header.Set(key, value)
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (l Login) Entitlement(asset AssetItem) (*Entitlement, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/entitlement/v2/asset/" + asset.DefaultPlayable.Id,
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Entitlement
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != nil {
      return nil, errors.New(result.Error.UserMessage)
   }
   return &result, nil
}

func (l Login) Playback(asset *AssetItem, title *Entitlement) (*Playback, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/playback/v1/preflight/asset/" + asset.DefaultPlayable.Id,
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("magine-play-deviceid", "!")
   req.Header.Set("magine-play-entitlementid", title.Token)
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   req.Header.Set("magine-play-devicemodel", "firefox 111.0 / windows 10")
   req.Header.Set("magine-play-deviceplatform", "firefox")
   req.Header.Set("magine-play-devicetype", "web")
   req.Header.Set("magine-play-drm", "widevine")
   req.Header.Set("magine-play-protocol", "dashs")
   // this value is important, with the wrong value you get random failures
   req.Header.Set("x-forwarded-for", "95.192.0.0")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
