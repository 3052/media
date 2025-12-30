package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)

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

func (d *Device) Register() error {
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
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(d)
}

func (a *AccountWithoutActiveProfile) SwitchProfile() (*Account, error) {
   data, err := json.Marshal(map[string]any{
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
   result := &Account{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
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
      "authorization", "Bearer "+d.Data.RegisterDevice.Token.AccessToken,
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &AccountWithoutActiveProfile{}
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

type AccountWithoutActiveProfile struct {
   Data struct {
      Login struct {
         Account struct {
            Profiles []struct {
               Id string
            }
         }
      }
   }
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // AccountWithoutActiveProfile
         }
      }
   }
}

type Device struct {
   Data struct {
      RegisterDevice struct {
         Token struct {
            AccessToken     string
            RefreshToken    string
            AccessTokenType string // Device
         }
      }
   }
}
