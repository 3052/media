package canal

import (
   "fmt"
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

func TestFields(t *testing.T) {
   var fields1 Fields
   err := fields1.New(film.stream)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", fields1.ObjectIds())
}
