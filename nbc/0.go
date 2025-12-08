package nbc

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
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
`

func Android(name int) (*Metadata, error) {
   data, err := json.Marshal(map[string]any{
      "query": graphql_compact(bonanza_page),
      "variables": map[string]any{
         "app":      "nbc",
         "name":     strconv.Itoa(name),
         "oneApp":   true,
         "platform": "android",
         "type":     "VIDEO",
         "userId":   "",
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://friendship.nbc.co/v2/graphql", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var body struct {
      Data struct {
         BonanzaPage struct {
            Metadata Metadata
         }
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&body)
   if err != nil {
      return nil, err
   }
   if err := body.Errors; len(err) >= 1 {
      return nil, errors.New(err[0].Message)
   }
   return &body.Data.BonanzaPage.Metadata, nil
}

type Metadata struct {
   MpxAccountId    int64 `json:",string"`
   MpxGuid         int64 `json:",string"`
   ProgrammingType string
}

func (m *Metadata) StreamInfo() (*StreamInfo, error) {
   req, _ := http.NewRequest("", "https://lemonade.nbc.com", nil)
   req.URL.Path = func() string {
      data := []byte("/v1/vod/")
      data = strconv.AppendInt(data, m.MpxAccountId, 10)
      data = append(data, '/')
      data = strconv.AppendInt(data, m.MpxGuid, 10)
      return string(data)
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
   info := &StreamInfo{}
   err = json.NewDecoder(resp.Body).Decode(info)
   if err != nil {
      return nil, err
   }
   return info, nil
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
      Host:   "drmproxy.digitalsvc.apps.nbcuni.com",
      Path:   "/drm-proxy/license/playready",
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

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func (s StreamInfo) Mpd() (*url.URL, []byte, error) {
   resp, err := http.Get(s.PlaybackUrl)
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

type StreamInfo struct {
   PlaybackUrl string // MPD
}
