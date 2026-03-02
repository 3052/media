package hboMax

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (l *Login) Fetch(st *http.Cookie) error {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host, // Refactored
      Path:   "/authentication/linkDevice/login",
   }
   req.Header = http.Header{}
   req.AddCookie(st)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

func FetchSt() (*http.Cookie, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     api_host, // Refactored
      Path:     "/token",
      RawQuery: "realm=bolt",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}
func FetchInitiate(st *http.Cookie, market string) (*Initiate, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   join("default.beam-", market, ".prd.api.discomax.com"),
      Path:   "/authentication/linkDevice/initiate",
   }
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.AddCookie(st)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Attributes Initiate
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.Attributes, nil
}
