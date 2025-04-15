package movistar

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

func (d *device) unmarshal(data Byte[device]) error {
   return json.Unmarshal(data, d)
}

type device string

// mullvad pass
func (t *token) device(oferta1 *oferta) (Byte[device], error) {
   req, err := http.NewRequest(
      "POST", "https://auth.dof6.com?qspVersion=ssp", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/accounts/")
      b.WriteString(oferta1.AccountNumber)
      b.WriteString("/devices/")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusCreated {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
