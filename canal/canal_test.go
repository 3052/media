package canal

import (
   "fmt"
   "testing"
   "time"
)

var tests = []string{
   "https://www.canalplus.cz/stream/film/argylle-tajny-agent/",
   "https://www.canalplus.cz/stream/series/mozart-v-dzungli/",
   "https://www.canalplus.cz/stream/series/silo/",
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      var fields1 fields
      err := fields1.New(test1)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", fields1.algolia_convert_tracking())
      time.Sleep(time.Second)
   }
}
