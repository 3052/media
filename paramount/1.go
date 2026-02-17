package paramount

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

// 1080p SL2000
// 1440p SL2000 + cookie
func PlayReady(at, contentId string, cookie *http.Cookie) (*SessionToken, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   if cookie != nil {
      req.AddCookie(cookie)
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/session-token.json"
   } else {
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json"
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

// 576p L3
func Widevine(at, contentId string) (*SessionToken, error) {
   var req http.Request
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.Path = "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}
