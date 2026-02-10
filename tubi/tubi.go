package tubi

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func (v *VideoResource) Dash() (*Dash, error) {
   resp, err := http.Get(v.Manifest.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (c *Content) Fetch(id int) error {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "uapi.adrise.tv",
      Path:   "/cms/content",
      RawQuery: url.Values{
         "content_id": {strconv.Itoa(id)},
         "deviceId":   {"!"},
         "platform":   {"android"},
         "video_resources[]": {
            "dash",
            "dash_widevine",
         },
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(c)
   if err != nil {
      return err
   }
   if len(c.VideoResources) == 0 {
      return errors.New("no video resources found")
   }
   return nil
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
