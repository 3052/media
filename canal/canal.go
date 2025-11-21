package canal

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
   "time"
)

var Transport = http.Transport{
   Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
   Proxy: func(req *http.Request) (*url.URL, error) {
      if path.Ext(req.URL.Path) == ".dash" {
         return nil, nil
      }
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
}

func (e *Episode) String() string {
   data := []byte("episode = ")
   data = strconv.AppendInt(data, e.Params.SeriesEpisode, 10)
   data = append(data, "\ntitle = "...)
   data = append(data, e.Title...)
   data = append(data, "\nid = "...)
   data = append(data, e.Id...)
   return string(data)
}

type Episode struct {
   Id string
   Params struct {
      SeriesEpisode int64
   }
   Title string
}

func (t *Ticket) Token(username, password string) (*Token, error) {
   value := map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cplogin.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   var client client_token
   err = client.New(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client.String())
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var tokenVar Token
   err = json.NewDecoder(resp.Body).Decode(&tokenVar)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, &tokenVar
   }
   return &tokenVar, nil
}

func (t *Token) Error() string {
   var data strings.Builder
   data.WriteString("label = ")
   data.WriteString(t.Label)
   data.WriteString("\nmessage = ")
   data.WriteString(t.Message)
   return data.String()
}

type Token struct {
   Label    string
   Message  string
   SsoToken string // this last one day
}

func FetchSession(ssoToken string) (SessionData, error) {
   data, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     ssoToken,
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://tvapi-hlm2.solocoo.tv/v1/session", "", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (t *Ticket) Fetch() error {
   data, err := json.Marshal(map[string]any{
      "deviceInfo": map[string]string{
         "brand":        "m7cp", // sg.ui.sso.fatal.internal_error
         "deviceModel":  "Firefox",
         "deviceOem":    "Firefox",
         "deviceSerial": device_serial,
         "deviceType":   "PC",
         "osVersion":    "Windows 10",
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cplogin.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   var client client_token
   err = client.New(req.URL, data)
   if err != nil {
      return err
   }
   req.Header.Set("authorization", client.String())
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(t)
   if err != nil {
      return err
   }
   if t.Message != "" {
      return errors.New(t.Message)
   }
   return nil
}

func (s *Session) Episodes(tracking_id string, season int64) ([]Episode, error) {
   req, _ := http.NewRequest("", "https://tvapi-hlm2.solocoo.tv/v1/assets", nil)
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.URL.RawQuery = func() string {
      data := []byte("limit=99&query=episodes,")
      data = append(data, tracking_id...)
      data = append(data, ",season,"...)
      data = strconv.AppendInt(data, season, 10)
      return string(data)
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Assets  []Episode
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if value.Message != "" {
      return nil, errors.New(value.Message)
   }
   return value.Assets, nil
}

type client_token struct {
   time int64
   sig  []byte
}

const (
   // clientKey is the public identifier for the client.
   clientKey = "web.NhFyz4KsZ54"
   // secretKey is the base64 encoded secret for HMAC.
   secretKey = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

func (c *client_token) New(address *url.URL, body []byte) error {
   bodyHash := sha256.Sum256(body)
   c.time = time.Now().Unix()
   decodedSecret, err := base64.RawURLEncoding.DecodeString(secretKey)
   if err != nil {
      return err
   }
   hasher := hmac.New(sha256.New, decodedSecret)
   fmt.Fprint(hasher, address)
   fmt.Fprint(hasher, base64.RawURLEncoding.EncodeToString(bodyHash[:]))
   fmt.Fprint(hasher, c.time)
   c.sig = hasher.Sum(nil)
   return nil
}

func (c *client_token) String() string {
   data := []byte("Client key=")
   data = append(data, clientKey...)
   data = append(data, ",time="...)
   data = fmt.Append(data, c.time)
   data = append(data, ",sig="...)
   data = base64.RawURLEncoding.AppendEncode(data, c.sig)
   return string(data)
}

func (s *Session) Player(tracking_id string) (*Player, error) {
   data, err := json.Marshal(map[string]any{
      "player": map[string]any{
         "capabilities": map[string]any{
            "drmSystems": []string{"Widevine"},
            "mediaTypes": []string{"DASH"},
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://tvapi-hlm2.solocoo.tv", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v1/assets/")
      data.WriteString(tracking_id)
      data.WriteString("/play")
      return data.String()
   }()
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var play Player
   err = json.NewDecoder(resp.Body).Decode(&play)
   if err != nil {
      return nil, err
   }
   if play.Message != "" {
      return nil, errors.New(play.Message)
   }
   return &play, nil
}

func TrackingId(address string) (string, error) {
   resp, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   const startKey = `data-algolia-convert-tracking="`
   _, after, found := strings.Cut(string(data), startKey)
   if !found {
      return "", fmt.Errorf("attribute key '%s' not found", startKey)
   }
   value, _, found := strings.Cut(after, `"`)
   if !found {
      return "", fmt.Errorf("could not find closing quote for the attribute")
   }
   return value, nil
}

type Player struct {
   Drm struct {
      LicenseUrl string
   }
   Message string
   Url     string // MPD
}

const device_serial = "!!!!"

func (p *Player) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Ticket struct {
   Message string
   Ticket  string
}

type Session struct {
   Message  string
   SsoToken string
   Token    string // this last one hour
}

func (s *Session) Unmarshal(data SessionData) error {
   err := json.Unmarshal(data, s)
   if err != nil {
      return err
   }
   if s.Message != "" {
      return errors.New(s.Message)
   }
   return nil
}

type SessionData []byte
