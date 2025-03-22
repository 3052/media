package main

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

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header["Accept"] = []string{"application/json, text/plain, */*"}
   req.Header["Accept-Language"] = []string{"en-US,en;q=0.5"}
   req.Header["Cache-Control"] = []string{"no-cache"}
   req.Header["Connection"] = []string{"keep-alive"}
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Host"] = []string{"m7cplogin.solocoo.tv"}
   req.Header["Origin"] = []string{"https://play.canalplus.cz"}
   req.Header["Pragma"] = []string{"no-cache"}
   req.Header["Referer"] = []string{"https://play.canalplus.cz/"}
   req.Header["Sec-Fetch-Dest"] = []string{"empty"}
   req.Header["Sec-Fetch-Mode"] = []string{"cors"}
   req.Header["Sec-Fetch-Site"] = []string{"cross-site"}
   req.Header["Te"] = []string{"trailers"}
   req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0"}
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "m7cplogin.solocoo.tv"
   req.URL.Path = "/login"
   req.URL.RawPath = ""
   value := url.Values{}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   var client1 client
   err := client1.New(req.URL, []byte(data))
   if err != nil {
      panic(err)
   }
   req.Header.Set("authorization", client1.String())
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

const data = `
{
 "provisionData": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJiciI6Im03Y3AiLCJ1cCI6Im03Y3AiLCJpYyI6dHJ1ZSwib3AiOiIxMDA0NiIsImRlIjoiYnJhbmRNYXBwaW5nIiwiaWF0IjoxNzQyNTIzNjU2LCJkcyI6IncxZjhhOGZiMC0wNWZiLTExZjAtYjVkYS1mMzJkMWNkNWRkZjcifQ.EOAV_4fk5cDsU8b1j80ni5F7N5Q7jhEfxhxvZ3oqOVM",
 "deviceInfo": {
  "osVersion": "Windows 10",
  "deviceModel": "Firefox",
  "deviceType": "PC",
  "deviceSerial": "w1f8a8fb0-05fb-11f0-b5da-f32d1cd5ddf7",
  "deviceOem": "Firefox",
  "devicePrettyName": "Firefox 128.0",
  "appVersion": "12.2",
  "language": "en_US",
  "brand": "m7cp",
  "country": "CZ"
 }
}
`
