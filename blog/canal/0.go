package canal

import (
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
)

func zero() (*http.Response, error) {
   const data = `
   {
      "deviceInfo": {
         "deviceModel": "Firefox",
         "deviceOem": "Firefox",
         "deviceSerial": "",
         "deviceType": "PC",
         "osVersion": "Windows 10"
      }
   }
   `
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "m7cplogin.solocoo.tv"
   req.URL.Path = "/login"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   var client1 client
   err := client1.New(req.URL, []byte(data))
   if err != nil {
      panic(err)
   }
   req.Header.Set("authorization", client1.String())
   return http.DefaultClient.Do(&req)
}

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
