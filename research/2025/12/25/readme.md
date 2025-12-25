# disney

here is what the web client does, note we can probably omit some of these calls.
first it does `registerDevice`:

~~~
POST https://disney.api.edge.bamgrid.com/graph/v1/device/graphql HTTP/2.0
authorization: Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu...

{
  "query": "mutation registerDevice($input: RegisterDeviceInput!) {\n      registerDevice(registerDevice: $input) {\n        grant {\n          grantType\n          assertion\n        },\n        token {\n          accessToken\n          accessTokenType\n          expiresIn\n          refreshToken\n          tokenType\n        },\n        session: activeSession {\n          sessionId\n          partnerName\n          device {\n            id\n            category\n            platform\n          }\n          profile {\n            id\n          }\n          experiments {\n            featureId\n            variantId\n            version\n          }\n          portabilityLocation {\n            countryCode\n            type\n          }\n          homeLocation {\n            adsSupported\n            countryCode\n          }\n          household {\n            householdScore\n          }\n          preferredMaturityRating {\n            impliedMaturityRating\n            ratingSystem\n          }\n          identity {\n            id\n          }\n          location {\n            adsSupported\n            type\n            countryCode\n            dma\n            asn\n            regionName\n            connectionType\n            zipCode\n          }\n        }\n      }\n    }",
  "variables": {
    "input": {
      "deviceProfile": "windows",
      "deviceFamily": "browser",
      "applicationRuntime": "firefox",
      "attributes": {
        "operatingSystem": "windows",
        "operatingSystemVersion": "10.0"
      }
    }
  }
}

HTTP/2.0 200 

{
  "data": {
    "registerDevice": {
      "token": {
        "accessToken": "...0Bel-WKWWmtysYXPlhssXbFal_9Lz7gykDZLCgdQWuckROvkJ...",
        "refreshToken": "...BzU_HikpzuPbDyTvaXpzDRmxS0n1NqR7e20tEjoJSfirpos-...",
        "accessTokenType": "Device"
      }
    }
  }
}
~~~
