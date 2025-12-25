package disney

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (r *register_device) login(email, password string) (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer " + r.Token.AccessToken)
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
   `, query_login, email, password)
   req.Body = io.NopCloser(strings.NewReader(data))
   return http.DefaultClient.Do(&req)
}

const query_login = `
mutation login($input: LoginInput!) {
  login(login: $input) {
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
