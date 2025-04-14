package movistar

import (
   "bytes"
   "encoding/json"
   "log"
   "net/http"
   "strings"
)

func (d device) session(init1 *init_data) (*http.Response, error) {
   data, err := json.Marshal(map[string]any{
      "contentID":  3427440,
      "drmMediaID": "1176568",
      "streamType": "AST",
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://alkasvaspub.imagenio.telefonica.net",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/asvas/ccs/")
      b.WriteString(init1.AccountNumber)
      b.WriteString("/SMARTTV_OTT/")
      b.WriteString(string(d))
      b.WriteString("/Session")
      return b.String()
   }()
   req.Header = http.Header{
      "content-type": {"application/json"},
      "x-hzid":       {init1.Token},
   }
   log.Print(req.Header)
   return http.DefaultClient.Do(req)
}
