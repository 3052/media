package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)

func (r *register_device) fetch() error {
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
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Add("Authorization", "Bearer " + client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(r)
}

type register_device struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken string
         }
      }
   }
}

const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

const query_register_device = `
mutation registerDevice($input: RegisterDeviceInput!) {
   registerDevice(registerDevice: $input) {
      token {
         accessToken
         refreshToken
      }
   }
}
`

