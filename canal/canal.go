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

func (f fields) object_ids() string {
   return f.get("objectIDs")
}

func (f fields) get(key string) string {
   for i, field := range f {
      if field == key {
         return f[i+1]
      }
   }
   return ""
}

type fields []string

func (f *fields) New(address string) error {
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
      return strings.ContainsRune(" :'", r)
   })
   return nil
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

type client struct {
   sig  []byte
   time int64
}

// us fail
// czech republic mullvad fail
// czech republic nord fail
// czech republic smart proxy pass
func (s session) play() (*play, error) {
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
   req.URL.Path = "/v1/assets/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi/play"
   req.Header = http.Header{
      "authorization": {"Bearer " + s.Token},
      "content-type":  {"application/json"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var play1 play
   err = json.NewDecoder(resp.Body).Decode(&play1)
   if err != nil {
      return nil, err
   }
   if play1.Message != "" {
      return nil, errors.New(play1.Message)
   }
   return &play1, nil
}

func (t *ticket) token(username, password string) (Byte[token], error) {
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
   return io.ReadAll(resp.Body)
}

type ticket struct {
   Message string
   Ticket  string
}

func (t *ticket) New() error {
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

func (t token) session() (*session, error) {
   data, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     t.SsoToken,
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
   session1 := &session{}
   err = json.NewDecoder(resp.Body).Decode(session1)
   if err != nil {
      return nil, err
   }
   return session1, nil
}

const device_serial = "!!!!"

type session struct {
   Token string
}

type play struct {
   Drm struct {
      LicenseUrl string
   }
   Message string
   Url     string
}

func (p *play) widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(p.Drm.LicenseUrl, "", bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

type token struct {
   SsoToken string
}

func (t *token) unmarshal(data Byte[token]) error {
   return json.Unmarshal(data, t)
}
