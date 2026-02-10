package sky

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

type service struct {
   s string
}

func (s service) Error() string {
   s.s = strings.TrimSuffix(s.s, ".")
   return strings.ToLower(s.s)
}

var not_available = service{
   "We're sorry our service is not available in your region yet.",
}

// x-forwarded-for fail
// mullvad.net fail
// proxy-seller.com pass
func sky_player(cookie *http.Cookie) ([]byte, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "show.sky.ch"
   req.URL.Path = "/de/SkyPlayerAjax/SkyPlayer"
   req.URL.Scheme = "https"
   values := url.Values{}
   values["id"] = []string{"2035"}
   values["contentType"] = []string{"2"}
   req.URL.RawQuery = values.Encode()
   req.Header["X-Requested-With"] = []string{"XMLHttpRequest"}
   req.AddCookie(cookie)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if strings.Contains(string(data), not_available.s) {
      return nil, not_available
   }
   return data, nil
}
