package max

import (
   "net/http"
   "strings"
)

func (d default_token) video() (*http.Response, error) {
   req, err := http.NewRequest(
      "", "https://default.any-amer.prd.api.discomax.com", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/content/videos/")
      b.WriteString("127b00c5-0131-4bac-b2d1-40762deefe09")
      b.WriteString("/activeVideoForShow")
      return b.String()
   }()
   req.Header.Set("authorization", "Bearer " + d.Data.Attributes.Token)
   return http.DefaultClient.Do(req)
}
