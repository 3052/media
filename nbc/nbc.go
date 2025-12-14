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

const drmProxySecret = "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"

// buildAuthQuery generates the signed query parameters (hash, time, device).
func buildAuthQuery(drmType string) string {
   timestamp := fmt.Sprint(time.Now().UnixMilli())
   mac := hmac.New(sha256.New, []byte(drmProxySecret))
   fmt.Fprint(mac, timestamp, drmType)
   hash := fmt.Sprintf("%x", mac.Sum(nil))
   return url.Values{
      "device": {"web"},
      "hash":   {hash},
      "time":   {timestamp},
   }.Encode()
}

func playReady() *url.URL {
   return &url.URL{
      Scheme:   "https",
      Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
      Path:     "/drm-proxy/license/playready",
      RawQuery: buildAuthQuery("playready"),
   }
}

func Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://drmproxy.digitalsvc.apps.nbcuni.com",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }

   req.URL.Path = "/drm-proxy/license/widevine"
   req.URL.RawQuery = buildAuthQuery("widevine")
   req.Header.Set("Content-Type", "application/octet-stream")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

func (s Stream) Mpd() (*url.URL, []byte, error) {
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

type Stream struct {
   PlaybackUrl string // MPD
}

// https://nbc.com/saturday-night-live/video/november-15-glen-powell/9000454161
func GetName(rawUrl string) (string, error) {
   parsed, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(parsed.Path, "/"), nil
}

func FetchMetadata(name string) (*Metadata, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_page,
      "variables": map[string]string{
         "app":      "nbc",
         "name":     name,
         "platform": "web",
         "type":     "VIDEO",
         "userId":   "",
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://friendship.nbc.com/v3/graphql", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Page struct {
            Metadata Metadata
         }
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data.Page.Metadata, nil
}

const query_page = `
query page(
   $app: NBCUBrands!
   $name: String!
   $platform: SupportedPlatforms!
   $type: PageType!
   $userId: String!
) {
  page(
    app: $app
    name: $name
    platform: $platform
    type: $type
    userId: $userId
  ) {
    metadata {
      ...on VideoPageMetaData {
        mpxAccountId
        mpxGuid
        programmingType
      }
    }
  }
}
`

type Metadata struct {
   MpxAccountId    int64 `json:",string"`
   MpxGuid         int64 `json:",string"`
   ProgrammingType string
}

func (m *Metadata) Stream() (*Stream, error) {
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
   result := &Stream{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
