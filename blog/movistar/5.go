package movistar

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strings"
)

func (d device) session(init1 *init_data, details1 *details) (*session, error) {
   data, err := json.Marshal(map[string]any{
      "contentID":  details1.Id,
      "drmMediaID": details1.VodItems[0].CasId,
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
      b.WriteByte('/')
      b.WriteString(device_type)
      b.WriteByte('/')
      b.WriteString(string(d))
      b.WriteString("/Session")
      return b.String()
   }()
   req.Header = http.Header{
      "content-type": {"application/json"},
      "x-hzid":       {init1.Token},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   session1 := &session{}
   err = json.NewDecoder(resp.Body).Decode(session1)
   if err != nil {
      return nil, err
   }
   return session1, nil
}

type session struct {
   ResultData struct {
      Ctoken string
   }
}
