package pluto

import (
   "strings"
   "testing"
)

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
