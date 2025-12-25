package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const query_register_device = `
mutation registerDevice($input: RegisterDeviceInput!) {
   registerDevice(registerDevice: $input) {
      token {
         accessToken
         refreshToken
         accessTokenType
      }
   }
}
`

type register_device struct {
   Token struct {
      AccessToken string
      RefreshToken string
      AccessTokenType string // Device
   }
}

func fetch_register_device() (*register_device, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_register_device,
      "variables": map[string]any{
         "input": map[string]any{
            "deviceProfile": "!",
            "deviceFamily": "!",
            "applicationRuntime": "!",
            "attributes": map[string]string{
               "operatingSystem": "",
               "operatingSystemVersion": "",
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Add("Authorization", "Bearer " + client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RegisterDevice register_device
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.RegisterDevice, nil
}

const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"
