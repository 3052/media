package main

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "play.binge.com.au"
   req.URL.Path = "/api/v3/play"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(body)
   req.Header["Authorization"] = []string{"Bearer eyJraWQiOiI3a0UxeCt4bE5xbFJabHNaMm9NeStQNnlBckU9IiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiJhdXRoMHw2N2NlMjU4MzY1Yjc1MjA4NzNlMmU0ZjIiLCJodHRwOi8vZm94c3BvcnRzLmNvbS5hdS9tYXJ0aWFuX2lkIjoiYXV0aDB8NjdjZTI1ODM2NWI3NTIwODczZTJlNGYyIiwiaHR0cHM6Ly9zdHJlYW1vdGlvbi5jb20uYXUvYWNjb3VudC9leHRlcm5hbC1pZGVudGl0aWVzIjp7fSwiaHR0cDovL2lyZGV0by5jb20vY29udHJvbC9qdGkiOiI1MDJlY2ExZC00YWNiLTQ2NmMtYWNjMC1iMjA3Y2I4MmJjMTgiLCJzZWNvbmRhcnlfa2V5IjoiNzMyZjhjYWE0NTMyMjE1NWUyOTgwNDEzODNiNDI5OTc3MGM3MjVkODk5MmMyMjc1ZTdjMjcyZWM4MzY2ZTNjOCIsImlzcyI6Imh0dHBzOi8vdG9rZW5zZXJ2aWNlLnN0cmVhbW90aW9uLmNvbS5hdS8iLCJodHRwOi8vaXJkZXRvLmNvbS9jb250cm9sL2FpZCI6ImZveHRlbG90dCIsImd0eSI6InBhc3N3b3JkIiwiYXVkIjoic3RyZWFtb3Rpb24uY29tLmF1IiwiaHR0cHM6Ly92aW1vbmQvZW50aXRsZW1lbnRzIjpbeyJzdHJlYW1jb3VudCI6MSwiYWRfc3VwcG9ydGVkIjp0cnVlLCJzdm9kIjoiMyIsInF1YWxpdHkiOiJmaGQifV0sImF6cCI6InBNODdUVVhLUXZTU3U5M3lkUmpEVHFCZ2RZZUNiZGhaIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBhZGRyZXNzIHBob25lIG9mZmxpbmVfYWNjZXNzIHVzZXI6cGhvbmVfdmVyaWZpZWQiLCJodHRwOi8vaXJkZXRvLmNvbS9jb250cm9sL2VudCI6W3siZXBpZCI6IldlYl9BUkVTXzFGSEQiLCJiaWQiOiJMSVRFIn1dLCJleHAiOjE3NDE2NjQzNzAsImlhdCI6MTc0MTY2NDA3MCwianRpIjoiZjBlNGQ1YjItNTFkYi00MTNiLThiZGQtNTY4MWRkNTIwZDM0IiwiaHR0cHM6Ly9hcmVzLmNvbS5hdS9zdGF0dXMiOnsidXBkYXRlZF9hdCI6IjIwMjUtMDMtMDlUMjM6Mzk6MzkuODMzWiIsInBwdl9ldmVudHMiOltdLCJhY2NvdW50X3N0YXR1cyI6IkFDVElWRV9TVUJTQ1JJUFRJT04iLCJzdWJfYWNjb3VudF9zdGF0dXMiOiJQQVlJTkdfU1VCU0NSSVBUSU9OIn19.X-lQPrIzRST7utJo92jg8eWMoqfErbz-juBH_0WDQhEbjSUALSZtlG1j1XXjYvBrZyq1zXJ8M_kePocLeweycoT9YQ1KOK-Q5A3bJoy-rvh16RDbr0KjU72bcp9ASaXVyRkzvMA2AtMzmqo6ebSPWBVqc89nROCPeF89lLyzDrYjLI7OTuNv_D29RS9DCgCMjDjQD2Tx9C3e1aZcbPpqnksfFVREGSWOvKTbf94rRUz3wMwWUEfLffnGsRvauANhOJmxluapVbwOuG9aKhaNNLYo4EvE8KgsFrcPkwcPpP3JVd2KMNCY7-T3iT4itsoP7LC_tLSbTnO72HeN_Gw0sw"}
   req.Header["Content-Type"] = []string{"application/json"}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

var body = strings.NewReader(`
{
   "assetId": "7738",
   "application": {
      "name": "binge"
   },
   "device": {
      "id": "50e785be-4c7f-4781-87e4-a3b4c75a3634"
   },
   "player": {
      "name": "VideoFS"
   }
}
`)
