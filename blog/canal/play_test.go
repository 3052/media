package canal

import (
   "41.neocities.org/widevine"
   "encoding/base64"
   "os"
   "testing"
)

var argylle = struct{
   key_id string
   url string
}{
   key_id: "8jU5F7LEqEP5pesDk/SaTw==",
   url: "play.canalplus.cz/player/d/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi",
}

func TestPlay(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/canal/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   session1, err := token1.session()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := session1.play()
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   key_id, err := base64.StdEncoding.DecodeString(argylle.key_id)
   if err != nil {
      t.Fatal(err)
   }
   pssh.KeyIds = [][]byte{key_id}
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = cdm.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = play1.widevine(data)
   if err != nil {
      t.Fatal(err)
   }
}
