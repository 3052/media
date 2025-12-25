package disney

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a *account_without_active_profile) switch_profile() (*account, error) {
   data, err := json.Marshal(map[string]any{
      "operationName": "switchProfile",
      "query": mutation_switch_profile,
      "variables": map[string]any{
         "input": map[string]string{
            "profileId": a.Data.Login.Account.Profiles[0].Id,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Extensions struct {
         Sdk struct {
            Token account
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Extensions.Sdk.Token, nil
}

type account_without_active_profile struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // AccountWithoutActiveProfile
         }
      }
   }
   Data struct {
      Login struct {
         Account struct {
            Profiles []struct {
               Id string
            }
         }
      }
   }
}

type account struct {
   AccessToken     string
   AccessTokenType string // Account
}

const mutation_switch_profile = `
mutation switchProfile($input: SwitchProfileInput!) {
    switchProfile(switchProfile: $input) {
      account {
        activeProfile {
          name
        }
      }
    }
  }
`

const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

type device struct {
   AccessToken     string
   RefreshToken    string
   AccessTokenType string // Device
}

const mutation_register_device = `
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

func register_device() (*device, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_register_device,
      "variables": map[string]any{
         "input": map[string]any{
            "deviceProfile":      "!",
            "deviceFamily":       "!",
            "applicationRuntime": "!",
            "attributes": map[string]string{
               "operatingSystem":        "",
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
   req.Header.Set("Authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RegisterDevice struct {
            Token device
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.RegisterDevice.Token, nil
}

func (d *device) login(email, password string) (*account_without_active_profile, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Header.Set("Authorization", "Bearer "+d.AccessToken)
   data := fmt.Sprintf(`
   {
     "operationName": "login",
     "query": %q,
     "variables": {
       "input": {
         "email": %q,
         "password": %q
       }
     }
   }
   `, mutation_login, email, password)
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &account_without_active_profile{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

const mutation_login = `
mutation login($input: LoginInput!) {
   login(login: $input) {
      account {
         profiles {
            id
         }
      }
   }
}
`
