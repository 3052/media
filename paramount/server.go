package paramount

import (
   "fmt"
   "net/http"
)

func GetAppSecret() (string, error) {
   const (
      targetURL  = "https://www.paramountplus.com"
      headerName = "x-real-server"
      usServer   = "us_www_web_prod_vip1"
      intlServer = "international_www_web_prod_vip1"
   )
   // 1. Perform a HEAD request.
   resp, err := http.Head(targetURL)
   if err != nil {
      return "", fmt.Errorf("failed to perform HEAD request to %s: %w", targetURL, err)
   }
   defer resp.Body.Close()
   // 2. Get the x-real-server response header.
   serverHeader := resp.Header.Get(headerName)
   // 3. Check the header value and return the corresponding string.
   switch serverHeader {
   case usServer:
      return AppSecrets[0].ComCbsApp, nil
   case intlServer:
      return AppSecrets[0].ComCbsCa, nil
   }
   // 4. Else, return an empty string and an error.
   return "", fmt.Errorf("unexpected value for header %q: got %q, want %q or %q",
      headerName, serverHeader, usServer, intlServer)
}
