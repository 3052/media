package canal

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strconv"
   "time"
)

const (
   key = "web.NhFyz4KsZ54"
   secret = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

type client struct {
   sig []byte
   unix int64
}

func (c *client) String() string {
   b := []byte("Client key=")
   b = append(b, key...)
   b = append(b, ",time="...)
   b = strconv.AppendInt(b, c.unix, 10)
   b = append(b, ",sig="...)
   b = base64.RawURLEncoding.AppendEncode(b, c.sig)
   return string(b)
}

func (c *client) New(ref *url.URL, body []byte) error {
   c.unix = time.Now().Unix()
   data := sha256.Sum256(body)
   secret1, err := base64.RawURLEncoding.DecodeString(secret)
   if err != nil {
      return err
   }
   hash := hmac.New(sha256.New, secret1)
   fmt.Fprint(hash, ref)
   fmt.Fprint(hash, base64.RawURLEncoding.EncodeToString(data[:]))
   fmt.Fprint(hash, c.unix)
   c.sig = hash.Sum(nil)
   return nil
}

type ticket struct {
   Ticket string
}

func (t *ticket) New() error {
   data, err := json.Marshal(map[string]any{
      "deviceInfo": map[string]string{
         "deviceModel": "Firefox",
         "deviceOem": "Firefox",
         "deviceSerial": "",
         "deviceType": "PC",
         "osVersion": "Windows 10",
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
   return json.NewDecoder(resp.Body).Decode(t)
}
