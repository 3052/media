package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const device_serial = "w1f8a8fb0-05fb-11f0-b5da-f32d1cd5ddf7"

func (t token) two() (*http.Response, error) {
   data, err := json.Marshal(map[string]string{
      "brand": "m7cp",
      "deviceSerial": device_serial,
      "deviceType": "PC",
      "ssoToken": t.SsoToken,
   })
   if err != nil {
      return nil, err
   }
   return http.Post(
      "https://tvapi-hlm2.solocoo.tv/v1/session", "", bytes.NewReader(data),
   )
}
