package draken

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (p *Playback) Dash() (*Dash, error) {
   resp, err := http.Get(p.Playlist)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

type Dash struct {
   Body []byte
   Url *url.URL
}

type Playback struct {
   Headers  map[string]string
   Playlist string // MPD
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

func (l *Login) Fetch(identity, accessKey string) error {
   data, err := json.Marshal(map[string]string{
      "accessKey": accessKey,
      "identity":  identity,
   })
   if err != nil {
      return err
   }
   resp, err := http.Post(
      "https://drakenfilm.se/api/auth/login", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

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
   magine_accesstoken.set(req.Header)
   x_forwarded_for.set(req.Header)
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

type AssetItem struct {
   DefaultPlayable struct {
      Id string
   }
}

var magine_accesstoken = header{
   "magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663",
}

var magine_play_devicemodel = header{
   "magine-play-devicemodel", "firefox 111.0 / windows 10",
}

var magine_play_deviceplatform = header{
   "magine-play-deviceplatform", "firefox",
}

var magine_play_devicetype = header{
   "magine-play-devicetype", "web",
}

var magine_play_drm = header{
   "magine-play-drm", "widevine",
}

var magine_play_protocol = header{
   "magine-play-protocol", "dashs",
}

// this value is important, with the wrong value you get random failures
var x_forwarded_for = header{
   "x-forwarded-for", "95.192.0.0",
}

func (h *header) set(head http.Header) {
   head.Set(h.key, h.value)
}

type header struct {
   key   string
   value string
}

type Entitlement struct {
   Error *struct {
      UserMessage string `json:"user_message"`
   }
   Token string
}

func (l Login) Entitlement(asset AssetItem) (*Entitlement, error) {
   var req http.Request
   req.Header = http.Header{}
   magine_accesstoken.set(req.Header)
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/entitlement/v2/asset/" + asset.DefaultPlayable.Id,
   }
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
   req.Header = http.Header{}
   x_forwarded_for.set(req.Header)
   magine_accesstoken.set(req.Header)
   magine_play_devicemodel.set(req.Header)
   magine_play_deviceplatform.set(req.Header)
   magine_play_devicetype.set(req.Header)
   magine_play_drm.set(req.Header)
   magine_play_protocol.set(req.Header)
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("magine-play-deviceid", "!")
   req.Header.Set("magine-play-entitlementid", title.Token)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/playback/v1/preflight/asset/" + asset.DefaultPlayable.Id,
   }
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

func (l Login) Widevine(play *Playback, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/playback/v1/widevine/license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   magine_accesstoken.set(req.Header)
   for key, value := range play.Headers {
      req.Header.Set(key, value)
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
