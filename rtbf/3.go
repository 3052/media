package rtbf

import (
   "bytes"
   "encoding/json"
   "net/http"
)

type GigyaLogin struct {
   SessionToken string
}

func (j *Jwt) Login() (*GigyaLogin, error) {
   data, err := json.Marshal(map[string]any{
      "device": map[string]string{
         "deviceId": "",
         "type":     "WEB",
      },
      "jwt": j.IdToken,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://exposure.api.redbee.live", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin"
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   gigya := &GigyaLogin{}
   err = json.NewDecoder(resp.Body).Decode(gigya)
   if err != nil {
      return nil, err
   }
   return gigya, nil
}
