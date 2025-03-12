package binge

import (
   "41.neocities.org/widevine"
   "os"
   "testing"
)

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 auth
   err = auth1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token, err := auth1.token()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := token.play()
   if err != nil {
      t.Fatal(err)
   }
   stream1, ok := play1.dash()
   if !ok {
      t.Fatal(".dash()")
   }
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   pssh.KeyIds = [][]byte{
      []byte(key_id),
   }
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = cdm.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = token.widevine(stream1, data)
   if err != nil {
      t.Fatal(err)
   }
}

const key_id = "\x1e\v\xdc\xfd\x069Nٝ\x9f\xe6#\xb0<\xa6\xd2"
