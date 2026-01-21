package paramount

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (s *SessionToken) Fetch(at, contentId string) error {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "www.paramountplus.com",
      Path: "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json",
      RawQuery: url.Values{
         "at":        {at},
         "contentId": {contentId},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return err
      }
      return errors.New(data.String())
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

func (s *SessionToken) PlayReady(at, contentId string) error {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "www.paramountplus.com",
      Path: "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json",
      RawQuery: url.Values{
         "at":        {at},
         "contentId": {contentId},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

// L3 540p
// L1 720p
func (s *SessionToken) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(string(data))
   }
   return data, nil
}
