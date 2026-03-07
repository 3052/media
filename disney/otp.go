package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
   "strconv"
)

func (d *Device) authenticate_with_otp(email string, passcode int) (*authenticate_with_otp, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_authenticate_with_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email": email,
            "passcode": strconv.Itoa(passcode),
         },
      },
      "operationName": "authenticateWithOtp",
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
         AuthenticateWithOtp authenticate_with_otp
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.AuthenticateWithOtp, nil
}

type authenticate_with_otp struct {
   ActionGrant string
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
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.RequestOtp, nil
}

const mutation_authenticate_with_otp = `
mutation authenticateWithOtp($input: AuthenticateWithOtpInput!) {
   authenticateWithOtp(authenticateWithOtp: $input) {
      actionGrant
      securityAction
      identity {
         personalInfo {
            dateOfBirth
            gender
         }
         flows {
            personalInfo {
               requiresCollection
               eligibleForCollection
            }
         }
      }
   }
}
`

const mutation_request_otp = `
mutation requestOtp($input: RequestOtpInput!) {
   requestOtp(requestOtp: $input) {
      accepted
   }
}
`
