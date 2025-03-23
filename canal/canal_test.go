package canal

import (
   "41.neocities.org/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "os/exec"
   "path"
   "strings"
   "testing"
)

func TestToken(t *testing.T) {
   data, err := exec.Command("password", "canalplus.cz").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var ticket1 ticket
   err = ticket1.New()
   if err != nil {
      t.Fatal(err)
   }
   data, err = ticket1.token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/canal/token", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}

func TestFields(t *testing.T) {
   var fields1 fields
   err := fields1.New(film.stream)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", fields1.object_ids())
}

var film = struct {
   key_id string
   player string
   stream string
}{
   key_id: "8jU5F7LEqEP5pesDk/SaTw==",
   player: "https://play.canalplus.cz/player/d/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi",
   stream: "https://www.canalplus.cz/stream/film/argylle-tajny-agent/",
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
   play1, err := session1.play(path.Base(film.player))
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
   key_id, err := base64.StdEncoding.DecodeString(film.key_id)
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
