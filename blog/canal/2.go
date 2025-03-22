package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const device_serial = "!!!!"

type token struct {
   Token string
}

func (s sso_token) token() (*token, error) {
   data, err := json.Marshal(map[string]string{
      "brand": "m7cp",
      "deviceSerial": device_serial,
      "deviceType": "PC",
      "ssoToken": s.SsoToken,
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://tvapi-hlm2.solocoo.tv/v1/session", "", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   token1 := &token{}
   err = json.NewDecoder(resp.Body).Decode(token1)
   if err != nil {
      return nil, err
   }
   return token1, nil
}
