package tubi

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

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
   req.Header.Set("vpn", "true")
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

type Content struct {
   Children     []*Content
   DetailedType string `json:"detailed_type"`
   Id           int    `json:",string"`
   SeriesId     int    `json:"series_id,string"`
   // these should already be in reverse order by resolution
   VideoResources []VideoResource `json:"video_resources"`
}

type Byte[T any] []byte

func (c *Content) Unmarshal(data Byte[Content]) error {
   return json.Unmarshal(data, c)
}
