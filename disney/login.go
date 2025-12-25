package disney

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a *account_without_active_profile) switch_profile() (*account, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer "+a.AccessToken)
   data := fmt.Sprintf(`
   {
     "query": %q,
     "variables": {
       "input": {
         "profileId": "ebb8f45f-fb18-4ebb-a2bf-fca32eb7fbb8"
       }
     },
     "operationName": "switchProfile"
   }
   `, mutation_switch_profile)
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
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

type account_without_active_profile struct {
   AccessToken     string
   AccessTokenType string // AccountWithoutActiveProfile
}

func (d *device) login(email, password string) (*account_without_active_profile, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer "+d.AccessToken)
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
   var result struct {
      Extensions struct {
         Sdk struct {
            Token account_without_active_profile
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Extensions.Sdk.Token, nil
}

const mutation_login = `
mutation login($input: LoginInput!) {
  login(login: $input) {
      actionGrant
  }
}
`
