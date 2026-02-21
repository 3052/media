package paramount

import (
   "errors"
   "net/http"
)

func FetchAppSecret() (string, error) {
   // 1. Perform a HEAD request.
   resp, err := http.Head("https://www.paramountplus.com")
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   switch resp.Header.Get("x-real-server") {
   case "us_www_web_prod_vip1":
      return AppSecrets[0].Us, nil
   case "international_www_web_prod_vip1":
      return AppSecrets[0].International, nil
   }
   return "", errors.New("unexpected or missing server header value")
}

var AppSecrets = []struct {
   Version       string
   Us            string
   International string
}{
   {
      Version:       "16.4.1",
      Us:            "7cd07f93a6e44cf7",
      International: "68b4475a49bed95a",
   },
   {
      Version:       "16.0.0",
      Us:            "9fc14cb03691c342",
      International: "6c68178445de8138",
   },
}
