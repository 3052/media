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

const mutation_login_with_action_grant = `
mutation loginWithActionGrant($input: LoginWithActionGrantInput!) {
   loginWithActionGrant(login: $input) {
      actionGrant
      account {
         activeProfile {
            id
         }
         profiles {
            id
            attributes {
               isDefault
               parentalControls {
                  isPinProtected
               }
            }
         }
      }
      activeSession {
         isSubscriber
      }
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

type authenticate_with_otp struct {
   ActionGrant string
}

func (d *Device) login_with_action_grant(action_grant string) (*http.Response, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_login_with_action_grant,
      "variables": map[string]any{
         "input": map[string]string{
            "actionGrant": action_grant,
         },
      },
      "operationName": "loginWithActionGrant",
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
   req.Header.Set("Authorization", "Bearer " + d.Token.AccessToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("Accept-Encoding", "identity")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Origin", "https://www.disneyplus.com")
   req.Header.Set("Referer", "https://www.disneyplus.com/")
   req.Header.Set("Sec-Fetch-Dest", "empty")
   req.Header.Set("Sec-Fetch-Mode", "cors")
   req.Header.Set("Sec-Fetch-Site", "cross-site")
   req.Header.Set("Te", "trailers")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0")
   req.Header.Set("X-Application-Version", "1.1.2")
   req.Header.Set("X-Bamsdk-Client-Id", "disney-svod-3d9324fc")
   req.Header.Set("X-Bamsdk-Platform", "javascript/windows/firefox")
   req.Header.Set("X-Bamsdk-Platform-Id", "browser")
   req.Header.Set("X-Bamsdk-Version", "34.4")
   req.Header.Set("X-Dss-Edge-Accept", "vnd.dss.edge+json; version=2")
   req.Header.Set("X-Request-Id", "85c72922-f5ff-4845-94a4-7584974d860a")
   req.Header.Set("X-Request-Yp-Id", "624b805dafc5c73635b1a216")
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
