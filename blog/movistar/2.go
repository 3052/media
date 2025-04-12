package movistar

import (
   "encoding/json"
   "errors"
   "net/http"
   "strings"
)

type device string

// mullvad pass
func (t *token) device(oferta1 *oferta) (device, error) {
   req, err := http.NewRequest(
      "POST", "https://auth.dof6.com?qspVersion=ssp", nil,
   )
   if err != nil {
      return "", err
   }
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/accounts/")
      b.WriteString(oferta1.AccountNumber)
      b.WriteString("/devices/")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusCreated {
      return "", errors.New(resp.Status)
   }
   var device1 device
   err = json.NewDecoder(resp.Body).Decode(&device1)
   if err != nil {
      return "", err
   }
   return device1, nil
}
