package disney

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

func register_device() (*http.Response, error) {
   data, err := json.Marshal(map[string]any{
      "operationName": "registerDevice",
      "query": query_register_device,
      "variables": map[string]any{
         "input": map[string]any{
            "deviceFamily": "browser",
            "applicationRuntime": "firefox",
            "deviceProfile": "windows",
            "deviceLanguage": "en",
            "devicePlatformId": "browser",
            "attributes": map[string]string{
               "brand": "web",
               "browserName": "firefox",
               "browserVersion": "128.0",
               "manufacturer": "n/a",
               "operatingSystem": "windows",
               "operatingSystemVersion": "10.0",
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
   req.Header.Add("Authorization", "Bearer " + client_api_key)
   return http.DefaultClient.Do(req)
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
