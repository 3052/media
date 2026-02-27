package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)


// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

// access token expires in 14400 seconds AKA 240 minutes AKA 4 hours. so using it
// properly we would:
// 1. refresh token
// 2. get movie
// 3. get movie
// 4. refresh token
// 5. get movie
// 6. get movie
// based on the duration that is idiotic and its simpler to just refresh every
// time


func (r *refresh_token) fetch() error {
   data, err := json.Marshal(map[string]any{
      "query": mutation_refresh_token,
      "variables": map[string]any{
         "input": map[string]string{
            "refreshToken": r.Extensions.Sdk.Token.RefreshToken,
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
   req.Header.Set("authorization", "Bearer " + client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(r)
}

type refresh_token struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken string
            RefreshToken string
         }
      }
   }
}

const mutation_refresh_token = `
mutation refreshToken($input: RefreshTokenInput!) {
   refreshToken(refreshToken: $input) {
      activeSession {
         sessionId
      }
   }
}
`
