package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const device_serial = "!!!!"

type session struct {
   Token string
}

func (t token) session() (*session, error) {
   data, err := json.Marshal(map[string]string{
      "brand": "m7cp",
      "deviceSerial": device_serial,
      "deviceType": "PC",
      "ssoToken": t.SsoToken,
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
   session1 := &session{}
   err = json.NewDecoder(resp.Body).Decode(session1)
   if err != nil {
      return nil, err
   }
   return session1, nil
}
