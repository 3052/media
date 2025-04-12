package movistar

import (
   "encoding/json"
   "errors"
   "net/http"
)

type device string

// XFF fail
// mullvad pass
func (t *token) device() (device, error) {
   req, err := http.NewRequest(
      "POST", "https://auth.dof6.com?qspVersion=ssp", nil,
   )
   if err != nil {
      return "", err
   }
   req.URL.Path = "/movistarplus/amazon.tv/accounts/00QSp000009M9gzMAC-L/devices/"
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
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
