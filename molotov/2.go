package molotov

import (
   "net/http"
   "net/url"
   "os"
)

func Two() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "fapi.molotov.tv"
   req.URL.Path = "/v2/channels/531/programs/15082/view"
   req.URL.Scheme = "https"
   req.Header.Set("x-molotov-agent", molotov_agent)
   value := url.Values{}
   value["access_token"] = []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMjgxODQxMDgiLCJhbGxvd2VkX2NpZHJzIjpbIjAuMC4wLjAvMCJdLCJleHBpcmVzIjoxNzQzMDM2MDgwLCJwcm9maWxlX2lkIjoiMjgxMzc5NjQiLCJzY29wZXMiOm51bGwsInVzZXJfaWQiOiIyODE4NDEwOCIsInYiOjF9.091390wNyt1_Mwbz9FhZtNYpNa6uASc8RJ1fwTb5fKE"}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}
