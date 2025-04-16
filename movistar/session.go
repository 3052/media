package movistar

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

type Session struct {
   ResultData struct {
      Ctoken string // ONE TIME USE
   }
   ResultText string
}

func (d Device) Session(init1 *InitData, details1 *Details) (*Session, error) {
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
   var value Session
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusCreated {
      return nil, errors.New(value.ResultText)
   }
   return &value, nil
}

func (s Session) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://wv-ottlic-f3.imagenio.telefonica.net",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/TFAESP/wvls/contentlicenseservice/v1/licenses"
   req.Header.Set("nv-authorizations", s.ResultData.Ctoken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
