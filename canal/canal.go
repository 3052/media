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
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

type Login struct {
   Label    string
   Message  string
   SsoToken string // this last one day
}

func (t *Ticket) Login(username, password string) (*Login, error) {
   data, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cp.login.solocoo.tv/login", bytes.NewReader(data),
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
   var result Login
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, &result
   }
   return &result, nil
}

func (l *Login) Error() string {
   var data strings.Builder
   data.WriteString("label = ")
   data.WriteString(l.Label)
   data.WriteString("\nmessage = ")
   data.WriteString(l.Message)
   return data.String()
}

func (p *Player) Mpd() (*url.URL, []byte, error) {
   resp, err := http.Get(p.Url)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

type Player struct {
   Drm struct {
      LicenseUrl string
   }
   Message string
   Subtitles []struct {
      Url string
   }
   Url     string // MPD
}

func (s *Session) Player(tracking string) (*Player, error) {
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
      data.WriteString(tracking)
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
   var result Player
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
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

func (s *Session) Episodes(tracking string, season int64) ([]Episode, error) {
   req, _ := http.NewRequest("", "https://tvapi-hlm2.solocoo.tv/v1/assets", nil)
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.URL.RawQuery = func() string {
      data := []byte("limit=99&query=episodes,")
      data = append(data, tracking...)
      data = append(data, ",season,"...)
      data = strconv.AppendInt(data, season, 10)
      return string(data)
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Assets  []Episode
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return result.Assets, nil
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

func (c *client_token) String() string {
   data := []byte("Client key=")
   data = append(data, clientKey...)
   data = append(data, ",time="...)
   data = fmt.Append(data, c.time)
   data = append(data, ",sig="...)
   data = base64.RawURLEncoding.AppendEncode(data, c.sig)
   return string(data)
}

const device_serial = "!!!!"

type Ticket struct {
   Message string
   Ticket  string
}

func (p *Player) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Session struct {
   Message  string
   SsoToken string
   Token    string // this last one hour
}

func (e *Episode) String() string {
   data := []byte("episode = ")
   data = strconv.AppendInt(data, e.Params.SeriesEpisode, 10)
   data = append(data, "\ntitle = "...)
   data = append(data, e.Title...)
   data = append(data, "\ndesc = "...)
   data = append(data, e.Desc...)
   data = append(data, "\ntracking = "...)
   data = append(data, e.Id...)
   return string(data)
}

type Episode struct {
   Desc string
   Id     string
   Params struct {
      SeriesEpisode int64
   }
   Title string
}

func (s *Session) Fetch(ssoToken string) error {
   data, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     ssoToken,
   })
   if err != nil {
      return err
   }
   resp, err := http.Post(
      "https://tvapi-hlm2.solocoo.tv/v1/session", "", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(s)
   if err != nil {
      return err
   }
   if s.Message != "" {
      return errors.New(s.Message)
   }
   return nil
}

func Tracking(address string) (string, error) {
   resp, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(resp.Status)
   }
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   const startKey = `data-algolia-convert-tracking="`
   _, after, found := strings.Cut(string(data), startKey)
   if !found {
      return "", fmt.Errorf("attribute key '%s' not found", startKey)
   }
   before, _, found := strings.Cut(after, `"`)
   if !found {
      return "", fmt.Errorf("could not find closing quote for the attribute")
   }
   return before, nil
}

func (c *client_token) New(address *url.URL, body []byte) error {
   hash := sha256.Sum256(body)
   c.time = time.Now().Unix()
   secret, err := base64.RawURLEncoding.DecodeString(secretKey)
   if err != nil {
      return err
   }
   hasher := hmac.New(sha256.New, secret)
   fmt.Fprint(hasher, address)
   fmt.Fprint(hasher, base64.RawURLEncoding.EncodeToString(hash[:]))
   fmt.Fprint(hasher, c.time)
   c.sig = hasher.Sum(nil)
   return nil
}
