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

// 1080p SL2000
// 1440p
func PlayReady(at, contentId string) (*SessionToken, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Path: "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json",
      Host: "www.paramountplus.com",
      RawQuery: url.Values{
         "at":        {at},
         "contentId": {contentId},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(data.String())
   }
   defer resp.Body.Close()
   var token SessionToken
   err = json.NewDecoder(resp.Body).Decode(&token)
   if err != nil {
      return nil, err
   }
   return &token, nil
}

// 540p L3
// 720p L1
func Widevine(at, contentId string) (*SessionToken, error) {
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
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(data.String())
   }
   defer resp.Body.Close()
   var token SessionToken
   err = json.NewDecoder(resp.Body).Decode(&token)
   if err != nil {
      return nil, err
   }
   return &token, nil
}

func (s *SessionToken) Send(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
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
