package movistar

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "strings"
)

// mullvad pass
func (o oferta) init_data(device1 device) (*init_data, error) {
   data, err := json.Marshal(map[string]string{
      "accountNumber": o.AccountNumber,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://clientservices.dof6.com?qspVersion=ssp",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/sdp/mediaPlayers/")
      b.WriteString(string(device1))
      b.WriteString("/initData")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   init1 := &init_data{}
   err = json.NewDecoder(resp.Body).Decode(init1)
   if err != nil {
      return nil, err
   }
   return init1, nil
}

type init_data struct {
   Token string
}
