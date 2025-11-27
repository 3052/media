package pluto

import (
   "fmt"
   "testing"
)

var tests = []string{
   "https://pluto.tv/on-demand/movies/6495eff09263a40013cf63a5",
   "https://pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6",
}

func TestPluto(t *testing.T) {
   t.Log(tests)
}
