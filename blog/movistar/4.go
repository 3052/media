package movistar

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func session() (*http.Response, error) {
   const data = `{"contentID":3427440,"drmMediaID":"1176568", "streamType":"AST"}`
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "alkasvaspub.imagenio.telefonica.net"
   req.URL.Path = "/asvas/ccs/00QSp000009M9gzMAC-L/SMARTTV_OTT/ea3585a776ed444d8677ad8be6ef0db3/Session"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   req.ContentLength = int64(len(data))
   req.Header["Content-Type"] = []string{"application/json"}
   req.Header["X-Hzid"] = []string{"eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2N2Y1Y2NlN2FkMDg3YjI1YzBmNjRhZGIiLCJpYXQiOjE3NDQ0MTIwNDQsImlzcyI6ImVhMzU4NWE3NzZlZDQ0NGQ4Njc3YWQ4YmU2ZWYwZGIzIiwiZXhwIjoxNzQ0NDU1MjQ0fQ.cYc7fzZFKT1CU5KWxuTZtEhy6CgP0rqFDBFdyjWwyJw"}
   return http.DefaultClient.Do(&req)
}
