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
   "strings"
   "time"
)

func (t Token) Session() (*Session, error) {
   value := map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     t.SsoToken,
   }
   data, err := json.MarshalIndent(value, "", " ")
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
   var session1 Session
   err = json.NewDecoder(resp.Body).Decode(&session1)
   if err != nil {
      return nil, err
   }
   if session1.Message != "" {
      return nil, errors.New(session1.Message)
   }
   return &session1, nil
}

func (t *Token) Unmarshal(data Byte[Token]) error {
   err := json.Unmarshal(data, t)
   if err != nil {
      return err
   }
   if t.Label != "sg.ui.sso.success.login" {
      return errors.New(t.Label)
   }
   return nil
}

type Token struct {
   Label string
   SsoToken string
}

func (t *Ticket) Token(username, password string) (Byte[Token], error) {
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
   return io.ReadAll(resp.Body)
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

type Session struct {
   Message string
   Token   string
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

type Byte[T any] []byte

func (p *Play) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Play struct {
   Drm struct {
      LicenseUrl string
   }
   Message string
   Url     string // MPD
}

const device_serial = "!!!!"

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
func (f Fields) ObjectIds() string {
   return f.get("objectIDs")
}

func (f *Fields) New(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   *f = strings.FieldsFunc(string(data), func(r rune) bool {
      return strings.ContainsRune(" ':[]", r)
   })
   return nil
}

func (f Fields) get(key string) string {
   for i, field := range f {
      if field == key {
         return f[i+1]
      }
   }
   return ""
}

type Fields []string

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

type client struct {
   sig  []byte
   time int64
}
