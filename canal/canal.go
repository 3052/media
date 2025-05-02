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
   "path"
   "strings"
   "time"
)

func (t *Ticket) Token(username, password string) (*Token, error) {
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
      "POST", "https://m7cplogin.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   var client1 client
   err = client1.New(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client1.String())
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var token1 Token
   err = json.NewDecoder(resp.Body).Decode(&token1)
   if err != nil {
      return nil, err
   }
   if token1.Label != "sg.ui.sso.success.login" {
      return nil, errors.New(token1.Label)
   }
   return &token1, nil
}

type Ticket struct {
   Message string
   Ticket  string
}

func (t *Ticket) New() error {
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
   var client1 client
   err = client1.New(req.URL, data)
   if err != nil {
      return err
   }
   req.Header.Set("authorization", client1.String())
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

type Token struct {
   Label    string
   SsoToken string // this last one day
}

type client struct {
   sig  []byte
   time int64
}

func (c *client) New(ref *url.URL, body []byte) error {
   body1 := sha256.Sum256(body)
   c.time = time.Now().Unix()
   secret1, err := base64.RawURLEncoding.DecodeString(secret)
   if err != nil {
      return err
   }
   hash := hmac.New(sha256.New, secret1)
   fmt.Fprint(hash, ref)
   fmt.Fprint(hash, base64.RawURLEncoding.EncodeToString(body1[:]))
   fmt.Fprint(hash, c.time)
   c.sig = hash.Sum(nil)
   return nil
}

func (c *client) String() string {
   b := []byte("Client key=")
   b = append(b, key...)
   b = append(b, ",time="...)
   b = fmt.Append(b, c.time)
   b = append(b, ",sig="...)
   b = base64.RawURLEncoding.AppendEncode(b, c.sig)
   return string(b)
}

const (
   key    = "web.NhFyz4KsZ54"
   secret = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

const device_serial = "!!!!"

type Byte[T any] []byte

func (p *Play) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Play) Unmarshal(data Byte[Play]) error {
   err := json.Unmarshal(data, p)
   if err != nil {
      return err
   }
   if p.Message != "" {
      return errors.New(p.Message)
   }
   return nil
}

type Play struct {
   Drm struct {
      LicenseUrl string
   }
   Message string
   Url     string // MPD
}

func NewSession(sso_token string) (Byte[Session], error) {
   data, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     sso_token,
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

func (s *Session) Unmarshal(data Byte[Session]) error {
   err := json.Unmarshal(data, s)
   if err != nil {
      return err
   }
   if s.Message != "" {
      return errors.New(s.Message)
   }
   return nil
}

type Session struct {
   Message  string
   SsoToken string
   Token    string // this last one hour
}

func (s *Session) Play(object_id string) (Byte[Play], error) {
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
      var b strings.Builder
      b.WriteString("/v1/assets/")
      b.WriteString(object_id)
      b.WriteString("/play")
      return b.String()
   }()
   // .Get .Set
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// https://play.canalplus.cz/player/d/Kc2fAJPVBKrayXNH2qQEuZV-94NggmNHxMQ0cpmT?
// parentId=SAVHw6HscpOmZ5tForujsLwVVWFKn8mobkGX5p2d
func ObjectId(data string) (string, error) {
   const prefix = "https://play.canalplus.cz/player/d/"
   if !strings.HasPrefix(data, prefix) {
      return "", fmt.Errorf("%q not found", prefix)
   }
   url2, err := url.Parse(data)
   if err != nil {
      return "", err
   }
   return path.Base(url2.Path), nil
}
