package movistar

import (
   "encoding/json"
   "errors"
   "net/http"
   "strings"
)

/*
request
GET https://auth.dof6.com/movistarplus/accounts/00QSp000009M9gzMAC-L/devices/fd242959339c49cf8c0a5054f653e49a?qspVersion=ssp HTTP/2.0
authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI3VTdlN3Y4QjhTOGg4bzlBIiwiYWNjb3VudE51bWJlciI6IjAwUVNwMDAwMDA5TTlnek1BQy1MIiwicm9sZSI6InVzZXIiLCJhcHIiOiJ3ZWJkYiIsImlzcyI6Imh0dHA6Ly93d3cubW92aXN0YXJwbHVzLmVzIiwiYXVkIjoiNDE0ZTE5MjdhMzg4NGY2OGFiYzc5ZjcyODM4MzdmZDEiLCJleHAiOjE3NDUzMzUyNDQsIm5iZiI6MTc0NDQ3MTI0NH0.6oraT4XCc5hXZpP4xkT0hCn3mtppQVwduo9NcRf01qw

response
{"Id":"fd242959339c49cf8c0a5054f653e49a","Name":"Amazon TV",
"DeviceTypeCode":"AMZTV","DeviceType":"Amazon TV",
"RegistrationDate":"2025-04-14T02:53:19.04Z","IsEnabled":true,
"IsInHomeZone":false,"IsInSsp":true,"IsPlaying":false,"ContentPlaying":null}
*/

type device string

// mullvad pass
func (t *token) device(oferta1 *oferta) (device, error) {
   req, err := http.NewRequest(
      "POST", "https://auth.dof6.com?qspVersion=ssp", nil,
   )
   if err != nil {
      return "", err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/accounts/")
      b.WriteString(oferta1.AccountNumber)
      b.WriteString("/devices/")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusCreated {
      return "", errors.New(resp.Status)
   }
   var device1 device
   err = json.NewDecoder(resp.Body).Decode(&device1)
   if err != nil {
      return "", err
   }
   return device1, nil
}
