package draken

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

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

///

func (n *Login) Entitlement(movie1 *Movie) (*Entitlement, error) {
   req, _ := http.NewRequest("POST", "https://client-api.magine.com", nil)
   req.URL.Path = "/api/entitlement/v2/asset/" + movie1.DefaultPlayable.Id
   req.Header.Set("authorization", "Bearer "+n.Token)
   magine_accesstoken.set(req.Header)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   title := &Entitlement{}
   err = json.NewDecoder(resp.Body).Decode(title)
   if err != nil {
      return nil, err
   }
   return title, nil
}

func (Login) Marshal(identity, key string) ([]byte, error) {
   data, err := json.Marshal(map[string]string{
      "accessKey": key,
      "identity":  identity,
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://drakenfilm.se/api/auth/login", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (m *Movie) New(custom_id string) error {
   data, err := json.Marshal(map[string]any{
      "query": graphql_compact(get_custom_id),
      "variables": map[string]string{
         "customId": custom_id,
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/apiql/v2",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   magine_accesstoken.set(req.Header)
   x_forwarded_for.set(req.Header)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Viewer struct {
            ViewableCustomId *Movie
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return err
   }
   if id := value.Data.Viewer.ViewableCustomId; id != nil {
      *m = *id
      return nil
   }
   return errors.New("ViewableCustomId")
}

type Login struct {
   Token string
}

type Entitlement struct {
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
            productionYear
            title
         }
      }
   }
}
` // do not do `query(`

type Movie struct {
   DefaultPlayable struct {
      Id string
   }
   ProductionYear int `json:",string"`
   Title          string
}

type Client struct {
   Login    *Login
   Playback *Playback
}

func (n *Login) Playback(movie1 *Movie, title *Entitlement) (*Playback, error) {
   req, _ := http.NewRequest("POST", "https://client-api.magine.com", nil)
   req.URL.Path = "/api/playback/v1/preflight/asset/" + movie1.DefaultPlayable.Id
   magine_accesstoken.set(req.Header)
   magine_play_devicemodel.set(req.Header)
   magine_play_deviceplatform.set(req.Header)
   magine_play_devicetype.set(req.Header)
   magine_play_drm.set(req.Header)
   magine_play_protocol.set(req.Header)
   req.Header.Set("authorization", "Bearer "+n.Token)
   req.Header.Set("magine-play-deviceid", "!")
   req.Header.Set("magine-play-entitlementid", title.Token)
   x_forwarded_for.set(req.Header)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   play := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

func (n *Login) Unmarshal(data []byte) error {
   return json.Unmarshal(data, n)
}

func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

type Playback struct {
   Headers  map[string]string
   Playlist string
}

func (p *Playback) Mpd() (*http.Response, error) {
   return http.Get(p.Playlist)
}

func (c *Client) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/playback/v1/widevine/license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   magine_accesstoken.set(req.Header)
   for key, value := range c.Playback.Headers {
      req.Header.Set(key, value)
   }
   req.Header.Set("authorization", "Bearer "+c.Login.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
