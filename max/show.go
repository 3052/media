package max

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
)

func show() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header["Authorization"] = []string{"Bearer eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1iMDM5YjUwYS1iYzMxLTRlNDAtOGUxYS1lYzJhOGFkOGMxZTYiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6YTFlOTNjZWUtZDQ0MC00ZTRkLWE0OGUtYzljNjRlNDg4YTIxIiwiaWF0IjoxNzQyNjU2NzkwLCJleHAiOjIwNTgwMTY3OTAsInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fZW1lYSIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6IjA3NzBmODU4LWRiZjQtNDM5NC1iNzdlLTMxMjhjZjA5ZTdiNyIsInZlcnNpb24iOiJ2MyIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiISJ9.qodQoDX2D4XD6MZ10aQl5i3FM3TlO6ijbAqW0oS7APC0J3UmcNVxsNBFse8FVa-5SUTPj4Eu83wyFq9YuBiRHQgJMsOggUhPbRhWtWuqQO6C23abrOY1yWFC0GRxAjP0HEPQWdOOt0CI6AFk5G1fzTDv8QAcU0tt5lfCJjo6nplJztMSXAoAu-yOhpNNyp1YF2c85CpNLrP36e4QzlH5oIeHCDEAxCvdh0Z5aP2bexOUscpVYR220Fd5qrRwWgbZuCrsXwRvPJbSzernIihsXZL1HlnDOmEebldp1IwvIeR2MHRmyHzsNomyMYmtNGJ9nHiAkMbYvrIuN_4GP2Vjog"}
   req.URL = &url.URL{}
   req.URL.Host = "default.prd.api.discomax.com"
   req.URL.Scheme = "https"
   req.URL.Path = "/cms/collections/227084608563650952176059252419027445293"
   req.URL.RawQuery = url.Values{
      "include":[]string{"default"},
      "pf[seasonNumber]":[]string{"1"},
      "pf[show.id]":[]string{"14f9834d-bc23-41a8-ab61-5c8abdbea505"},
   }.Encode()
   return http.DefaultClient.Do(&req)
}
