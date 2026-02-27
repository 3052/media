// to change location you must log in again
package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strings"
)

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name = ")
   data.WriteString(p.Name)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

type Profile struct {
   Name string
   Id   string
}

type AccountWithoutActiveProfile struct {
   Data struct {
      Login struct {
         Account struct {
            Profiles []struct {
               Id   string
               Name string
            }
         }
      }
   }
   Errors     []Error
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // AccountWithoutActiveProfile
         }
      }
   }
}

func (a *AccountWithoutActiveProfile) SwitchProfile(profileId string) (*Account, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_switch_profile,
      "variables": map[string]any{
         "input": map[string]string{
            "profileId": profileId,
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
   result := &Account{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
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

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

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

const mutation_login = `
mutation login($input: LoginInput!) {
   login(login: $input) {
      account {
         profiles {
            id
            name
         }
      }
   }
}
`

type Account struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // Account
         }
      }
   }
}

func (d *Device) Login(email, password string) (*AccountWithoutActiveProfile, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_login,
      "variables": map[string]any{
         "input": map[string]string{
            "email":    email,
            "password": password,
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
   req.Header.Set(
      "authorization", "Bearer "+d.Token.AccessToken,
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result AccountWithoutActiveProfile
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

type Device struct {
   Token struct {
      AccessToken     string
      RefreshToken    string
      AccessTokenType string // Device
   }
}

func RegisterDevice() (*Device, error) {
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
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RegisterDevice Device
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.RegisterDevice, nil
}
