package spotify

import (
   "154.pages.dev/protobuf"
   "bytes"
   "errors"
   "io"
   "net/http"
)

// const client_id = "9a8d2f0ce77a4e248bb71fefcb557637"
const client_id = "58bd3c95768941ea9eb4350aaa033eb3"

func (r *login_response) New(username, password string) error {
   if username == "" {
      return errors.New("username")
   }
   if password == "" {
      return errors.New("password")
   }
   var m protobuf.Message
   m.Add(1, func(m *protobuf.Message) {
      m.AddBytes(1, []byte(client_id))
   })
   m.Add(101, func(m *protobuf.Message) {
      m.AddBytes(1, []byte(username))
      m.AddBytes(2, []byte(password))
   })
   req, err := http.NewRequest(
      "POST", "https://login5.spotify.com/v3/login", bytes.NewReader(m.Encode()),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/x-protobuf")
   req.Header.Set("user-agent", "Symfony HttpClient (Curl)")
   req.Header.Set("accept", "*/*")
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   return r.m.Consume(data)
}

type login_response struct {
   m protobuf.Message
}
