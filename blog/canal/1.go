package canal

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
)

func (s *sso_token) unmarshal(data Byte[sso_token]) error {
   return json.Unmarshal(data, s)
}

type sso_token struct {
   SsoToken string
}

type Byte[T any] []byte

func (t ticket) token(username, password string) (Byte[sso_token], error) {
   data, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cplogin.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   var client1 client
   err = client1.New(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client1.String())
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

