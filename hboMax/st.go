package hboMax

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

type St struct {
   Cookie *http.Cookie
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (s St) Login() (*Login, error) {
   var req http.Request
   req.Header = http.Header{}
   req.AddCookie(s.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host, // Refactored
      Path:   "/authentication/linkDevice/login",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Login{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (s *St) Fetch() error {
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
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         s.Cookie = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}

func (s St) Initiate(market string) (*Initiate, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.AddCookie(s.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   join("default.beam-", market, ".prd.api.discomax.com"),
      Path:   "/authentication/linkDevice/initiate",
   }
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
