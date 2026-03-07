package disney

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func authenticate_with_otp() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer " + bearer)
   const data = `
   {
     "query": "\nmutation authenticateWithOtp($input: AuthenticateWithOtpInput!) {\n  authenticateWithOtp(authenticateWithOtp: $input) {\n      actionGrant\n      securityAction\n      identity {\n        personalInfo {\n          dateOfBirth\n          gender\n        }\n        flows {\n          personalInfo {\n            requiresCollection \n            eligibleForCollection\n          }\n        }\n      }\n    }\n  }\n",
     "variables": {
       "input": {
         "email": "27@riseup.net",
         "passcode": "839727"
       }
     },
     "operationName": "authenticateWithOtp"
   }
   `
   req.Body = io.NopCloser(strings.NewReader(data))
   return http.DefaultClient.Do(&req)
}

const mutation_request_otp = `
mutation requestOtp($input: RequestOtpInput!) {
   requestOtp(requestOtp: $input) {
      accepted
   }
}
`

func (d *Device) RequestOtp(email string) error {
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
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer " + device.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   return err
}

