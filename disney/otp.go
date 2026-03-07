package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)

type InactiveAccount struct {
   Data struct {
      Login *struct {
         Account struct {
            Profiles []Profile
         }
      }
      LoginWithActionGrant *struct {
         Account struct {
            Profiles []Profile
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

// passcode can start with 0
func (d *Device) AuthenticateWithOtp(email, passcode string) (*AuthenticateWithOtp, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_authenticate_with_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email": email,
            "passcode": passcode,
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
   req.Header.Set("authorization", "Bearer " + d.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         AuthenticateWithOtp AuthenticateWithOtp
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
   return &result.Data.AuthenticateWithOtp, nil
}

func (d *Device) LoginWithActionGrant(actionGrant string) (*InactiveAccount, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_login_with_action_grant,
      "variables": map[string]any{
         "input": map[string]string{
            "actionGrant": actionGrant,
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
   req.Header.Set("authorization", "Bearer " + d.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &InactiveAccount{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
func (d *Device) RequestOtp(email string) (*RequestOtp, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_request_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email": email,
            "reason": "Login",
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
   req.Header.Set("authorization", "Bearer " + d.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RequestOtp RequestOtp
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
   return &result.Data.RequestOtp, nil
}

func (r RequestOtp) String() string {
   if r.Accepted {
      return "accepted = true"
   }
   return "accepted = false"
}

type RequestOtp struct {
   Accepted bool
}

const mutation_request_otp = `
mutation requestOtp($input: RequestOtpInput!) {
   requestOtp(requestOtp: $input) {
      accepted
   }
}
`

type AuthenticateWithOtp struct {
   ActionGrant string
}

const mutation_authenticate_with_otp = `
mutation authenticateWithOtp($input: AuthenticateWithOtpInput!) {
   authenticateWithOtp(authenticateWithOtp: $input) {
      actionGrant
   }
}
`

const mutation_login_with_action_grant = `
mutation loginWithActionGrant($input: LoginWithActionGrantInput!) {
   loginWithActionGrant(login: $input) {
      account {
         profiles {
            id
            name
         }
      }
   }
}
`
