package movistar

import (
   "encoding/json"
   "errors"
   "net/http"
)

type oferta struct {
   AccountNumber string
}

// mullvad pass
func (t *token) oferta() (*oferta, error) {
   req, _ := http.NewRequest("", "https://auth.dof6.com", nil)
   req.URL.Path = "/movistarplus/api/devices/amazon.tv/users/authenticate"
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Ofertas []oferta
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Ofertas[0], nil
}
