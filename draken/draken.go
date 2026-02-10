package draken

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func FetchMovie(customId string) (*MovieItem, error) {
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
            ViewableCustomId *MovieItem
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
   Token string
   Error *Error
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("message = ")
   data.WriteString(e.Message)
   data.WriteString("\nuser message = ")
   data.WriteString(e.UserMessage)
   return data.String()
}

type Error struct {
   Message string
   UserMessage string `json:"user_message"`
}

type Login struct {
   Message string
   Token string
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

type Dash struct {
   Body []byte
   Url *url.URL
}

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

type MovieItem struct {
   DefaultPlayable struct {
      Id string
   }
}

func (l *Login) Widevine(play *Playback, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/playback/v1/widevine/license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   req.Header.Set("magine-play-deviceid", "!")
   req.Header.Set("magine-play-devicemodel", "firefox 111.0 / windows 10")
   req.Header.Set("magine-play-deviceplatform", "firefox")
   req.Header.Set("magine-play-devicetype", "web")
   req.Header.Set("magine-play-drm", "widevine")
   req.Header.Set("magine-play-protocol", "dashs")
   req.Header.Set("magine-play-session", play.Headers.MaginePlaySession)
   req.Header.Set(
      "magine-play-entitlementId", play.Headers.MaginePlayEntitlementId,
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Playback struct {
   Headers struct {
      MaginePlayEntitlementId string `json:"Magine-Play-EntitlementId"`
      MaginePlaySession string `json:"Magine-Play-Session"`
   }
   Playlist string // MPD
}

func (l *Login) Playback(movie *MovieItem, title *Entitlement) (*Playback, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/playback/v1/preflight/asset/" + movie.DefaultPlayable.Id,
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

func (l *Login) Fetch(identity, accessKey string) error {
   data, err := json.Marshal(map[string]string{
      "accessKey": accessKey,
      "identity":  identity,
   })
   if err != nil {
      return err
   }
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/login/v2/auth/email",
   }
   req.Header = http.Header{}
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(l)
   if err != nil {
      return err
   }
   if l.Message != "" {
      return errors.New(l.Message)
   }
   return nil
}

func (l *Login) Entitlement(movie *MovieItem) (*Entitlement, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "client-api.magine.com",
      Path: "/api/entitlement/v2/asset/" + movie.DefaultPlayable.Id,
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
      return nil, result.Error
   }
   return &result, nil
}
