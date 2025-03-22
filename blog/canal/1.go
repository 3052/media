package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

type token struct {
   SsoToken string
}

func (t ticket) token(username, password string) (*token, error) {
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
   token1 := &token{}
   err = json.NewDecoder(resp.Body).Decode(token1)
   if err != nil {
      return nil, err
   }
   return token1, nil
}
