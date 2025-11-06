package mubi

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
)

var ClientCountry = "US"

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}

type hello struct {
   Token string
   User  struct {
      Id int
   }
}

func (l *LinkCode) hello() (helloData, error) {
   data, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.mubi.com/v3/authenticate", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type helloData []byte

func (h *hello) Unmarshal(data helloData) error {
   return json.Unmarshal(data, h)
}
