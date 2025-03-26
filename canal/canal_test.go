package canal

import (
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

var film = struct {
   key_id string
   player string
   stream string
}{
   key_id: "8jU5F7LEqEP5pesDk/SaTw==",
   player: "https://play.canalplus.cz/player/d/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi",
   stream: "https://www.canalplus.cz/stream/film/argylle-tajny-agent/",
}

func TestToken(t *testing.T) {
   data, err := exec.Command("password", "canalplus.cz").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var ticket1 Ticket
   err = ticket1.New()
   if err != nil {
      t.Fatal(err)
   }
   data, err = ticket1.Token(username, password)
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
   var fields1 Fields
   err := fields1.New(film.stream)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", fields1.ObjectIds())
}
