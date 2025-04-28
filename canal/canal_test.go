package canal

import (
   "fmt"
   "testing"
)

/*
https://www.canalplus.cz/stream/series/silo/
*/

var film = struct {
   player string
   stream string
}{
   player: "https://play.canalplus.cz/player/d/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi",
   stream: "https://www.canalplus.cz/stream/film/argylle-tajny-agent/",
}

func Test(t *testing.T) {
   fmt.Println(film)
}
