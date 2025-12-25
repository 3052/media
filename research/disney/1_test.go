package disney

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   var play playback
   err := play.fetch()
   if err != nil {
      t.Fatal(err)
   }
   for i, source := range play.Stream.Sources {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Printf("%+v\n", source)
   }
}
