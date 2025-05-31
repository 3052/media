package rakuten

import (
   "fmt"
   "os"
   "testing"
)

func TestEpisode(t *testing.T) {
   data, err := os.ReadFile("address")
   if err != nil {
      t.Fatal(err)
   }
   var web address
   err = web.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   episodes, err := web.episodes(test.season)
   if err != nil {
      t.Fatal(err)
   }
   for _, episode1 := range episodes {
      fmt.Printf("%+v\n", episode1)
   }
}
