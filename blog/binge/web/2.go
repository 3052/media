package main

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

// residential proxy
func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "tokenservice.streamotion.com.au"
   req.URL.Path = "/oauth/token"
   req.URL.Scheme = "https"
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["Authorization"] = []string{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik56aEJPVFJHT0RjNE1FUkRSRFJEUTBVd1FrVkdNRGt4TVVVNVF6RTRRa0UzTkVVMk1rVkRRZyJ9.eyJzZWNvbmRhcnlfa2V5IjoiNzMyZjhjYWE0NTMyMjE1NWUyOTgwNDEzODNiNDI5OTc3MGM3MjVkODk5MmMyMjc1ZTdjMjcyZWM4MzY2ZTNjOCIsImlzcyI6Imh0dHBzOi8vYXV0aC5zdHJlYW1vdGlvbi5jb20uYXUvIiwic3ViIjoiYXV0aDB8NjdjZTI1ODM2NWI3NTIwODczZTJlNGYyIiwiYXVkIjpbInN0cmVhbW90aW9uLmNvbS5hdSIsImh0dHBzOi8vcHJvZC1tYXJ0aWFuLmZveHNwb3J0cy0xYi1wcm9kLmF1dGgwYXBwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3NDE1NzE1NDgsImV4cCI6MTc0MTU5MzE0OCwic2NvcGUiOiJvcGVuaWQgZW1haWwgZHJtOmxvdyBvZmZsaW5lX2FjY2VzcyB1c2VyOnBob25lX3ZlcmlmaWVkIiwiYXpwIjoicE04N1RVWEtRdlNTdTkzeWRSakRUcUJnZFllQ2JkaFoifQ.XcEfhjhu5Bwkm-d6Bg-Z3sTDWNr4wQFkt6ns-_lPbaoE6SUHGO8CNmLxK4m-vCdnGus4_bXlKevMnYohZhDGMiQwW-XbCg3FyCWAp-8K-3cMkZ49-AL4YxHZAwZE5HaQduqALSjQbqOE3-PKlpK7hY1Bf1W0qM0InSdV-DzdcxfDUsPBDcFca50uzSyPo-TPSDEvqlOviLTOgHlqByjV7sWArkpeXZGZZRkuAnmDB1MeHp0Z_CgV0OYitcU5zcBMnVyoEDuABlWZGcwqc1G3kiLlH6F-zO6BVRhXZLLrC_4Jz3sUJWhVe4nWPOdhQLer4wmD_Oq8isgbm9BAObPKCA"}
   req.Body = io.NopCloser(body)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

var body = strings.NewReader(`
{
   "client_id": "pM87TUXKQvSSu93ydRjDTqBgdYeCbdhZ"
}
`)
