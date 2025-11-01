package nbc

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
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
   Proxy: func(req *http.Request) (*url.URL, error) {
      if path.Ext(req.URL.Path) != ".mp4" {
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
}

func playReady() *url.URL {
   now := fmt.Sprint(time.Now().UnixMilli())
   hash := func() string {
      secret := hmac.New(sha256.New, []byte(drm_proxy_secret))
      fmt.Fprint(secret, now, "playready")
      return fmt.Sprintf("%x", secret.Sum(nil))
   }()
   return &url.URL{
      Scheme: "https",
      Host: "drmproxy.digitalsvc.apps.nbcuni.com",
      Path: "/drm-proxy/license/playready",
      RawQuery: url.Values{
         "device": {"web"},
         "hash":   {hash},
         "time":   {now},
      }.Encode(),
   }
}

func Widevine(data []byte) ([]byte, error) {
   now := fmt.Sprint(time.Now().UnixMilli())
   hash := func() string {
      hash1 := hmac.New(sha256.New, []byte(drm_proxy_secret))
      fmt.Fprint(hash1, now, "widevine")
      return fmt.Sprintf("%x", hash1.Sum(nil))
   }()
   req, err := http.NewRequest(
      "POST", "https://drmproxy.digitalsvc.apps.nbcuni.com",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/drm-proxy/license/widevine"
   req.URL.RawQuery = url.Values{
      "device": {"web"},
      "hash":   {hash},
      "time":   {now},
   }.Encode()
   req.Header.Set("content-type", "application/octet-stream")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

const drm_proxy_secret = "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"

type Vod struct {
   PlaybackUrl string // MPD
}

func (m *Metadata) New(guid int) error {
   value := map[string]any{
      "query": graphql_compact(bonanza_page),
      "variables": map[string]any{
         "app": "nbc",
         "name": strconv.Itoa(guid),
         "oneApp": true,
         "platform": "android",
         "type": "VIDEO",
         "userId": "",
      },
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return err
   }
   resp, err := http.Post(
      "https://friendship.nbc.co/v2/graphql", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var value1 struct {
      Data struct {
         BonanzaPage struct {
            Metadata Metadata
         }
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value1)
   if err != nil {
      return err
   }
   if err := value1.Errors; len(err) >= 1 {
      return errors.New(err[0].Message)
   }
   *m = value1.Data.BonanzaPage.Metadata
   return nil
}

const bonanza_page = `
query bonanzaPage(
   $app: NBCUBrands!
   $name: String!
   $oneApp: Boolean
   $platform: SupportedPlatforms!
   $type: EntityPageType!
   $userId: String!
) {
   bonanzaPage(
      app: $app
      name: $name
      oneApp: $oneApp
      platform: $platform
      type: $type
      userId: $userId
   ) {
      metadata {
         ... on VideoPageData {
            mpxAccountId
            mpxGuid
            programmingType
         }
      }
   }
}
` // do not use `query(`

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func (m *Metadata) Vod() (*Vod, error) {
   req, _ := http.NewRequest("", "https://lemonade.nbc.com", nil)
   req.URL.Path = func() string {
      b := []byte("/v1/vod/")
      b = strconv.AppendInt(b, m.MpxAccountId, 10)
      b = append(b, '/')
      b = strconv.AppendInt(b, m.MpxGuid, 10)
      return string(b)
   }()
   req.URL.RawQuery = url.Values{
      "platform":        {"web"},
      "programmingType": {m.ProgrammingType},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   video := &Vod{}
   err = json.NewDecoder(resp.Body).Decode(video)
   if err != nil {
      return nil, err
   }
   return video, nil
}

type Metadata struct {
   MpxAccountId     int64 `json:",string"`
   MpxGuid          int64 `json:",string"`
   ProgrammingType  string
}
