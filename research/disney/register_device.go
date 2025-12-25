package disney

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func register_device() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/graph/v1/device/graphql"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84")
   data := fmt.Sprintf(`
   {
     "operationName": "registerDevice",
     "query": %q,
     "variables": {
       "input": {
         "deviceFamily": "browser",
         "applicationRuntime": "firefox",
         "deviceProfile": "windows",
         "deviceLanguage": "en",
         "devicePlatformId": "browser",
         "attributes": {
           "brand": "web",
           "browserName": "firefox",
           "browserVersion": "128.0",
           "manufacturer": "n/a",
           "operatingSystem": "windows",
           "operatingSystemVersion": "10.0"
         }
       }
     }
   }
   `, query_register_device)
   req.Body = io.NopCloser(strings.NewReader(data))
   return http.DefaultClient.Do(&req)
}

const query_register_device = `
mutation registerDevice($input: RegisterDeviceInput!) {
      registerDevice(registerDevice: $input) {
        grant {
          grantType
          assertion
        },
        token {
          accessToken
          accessTokenType
          expiresIn
          refreshToken
          tokenType
        },
        session: activeSession {
          sessionId
          partnerName
          device {
            id
            category
            platform
          }
          profile {
            id
          }
          experiments {
            featureId
            variantId
            version
          }
          portabilityLocation {
            countryCode
            type
          }
          homeLocation {
            adsSupported
            countryCode
          }
          household {
            householdScore
          }
          preferredMaturityRating {
            impliedMaturityRating
            ratingSystem
          }
          identity {
            id
          }
          location {
            adsSupported
            type
            countryCode
            dma
            asn
            regionName
            connectionType
            zipCode
          }
        }
      }
    }
`
