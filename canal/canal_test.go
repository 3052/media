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

func Test(t *testing.T) {
   fmt.Println(film)
}
