package pluto

import (
   "strings"
   "testing"
)
package pluto

import "testing"

var tests = []string{
   "https://pluto.tv/on-demand/movies/6495eff09263a40013cf63a5",
   "https://pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6",
}

func TestPluto(t *testing.T) {
   t.Log(tests)
}
func Test(t *testing.T) {
   var value Series
   err := value.Fetch("6495eff09263a40013cf63a5")
   if err != nil {
      t.Fatal(err)
   }
   _, data, err := value.Mpd()
   if err != nil {
      t.Fatal(err)
   }
   const height = ` height="1080"`
   if !strings.Contains(string(data), height) {
      t.Fatal(height)
   }
}
