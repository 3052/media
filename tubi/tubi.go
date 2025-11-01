package tubi

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
)

var Transport = http.Transport{
   Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
   Proxy: func(req *http.Request) (*url.URL, error) {
      if path.Ext(req.URL.Path) != ".mp4" {
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
}

func (c *Content) Unmarshal(data Byte[Content]) error {
   err := json.Unmarshal(data, c)
   if err != nil {
      return err
   }
   if len(c.VideoResources) == 0 {
      return errors.New("video_resources")
   }
   return nil
}

type Content struct {
   Children     []*Content
   DetailedType string `json:"detailed_type"`
   Id           int    `json:",string"`
   SeriesId     int    `json:"series_id,string"`
   // these should already be in reverse order by resolution
   VideoResources []VideoResource `json:"video_resources"`
}

func NewContent(id int) (Byte[Content], error) {
   req, _ := http.NewRequest("", "https://uapi.adrise.tv/cms/content", nil)
   req.URL.RawQuery = url.Values{
      "content_id": {strconv.Itoa(id)},
      "deviceId":   {"!"},
      "platform":   {"android"},
      "video_resources[]": {
         "dash",
         "dash_widevine",
      },
   }.Encode()
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (v *VideoResource) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      v.LicenseServer.Url, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type VideoResource struct {
   LicenseServer *struct {
      Url string
   } `json:"license_server"`
   Manifest struct {
      Url string // MPD
   }
   Type string
}

type Byte[T any] []byte
