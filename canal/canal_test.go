package canal

import (
   "fmt"
   "testing"
   "time"
)

func TestAssets(t *testing.T) {
   for _, test1 := range tests {
      fmt.Println(test1.url)
      assets1, err := assets(test1.id, 1)
      if err != nil {
         t.Fatal(err)
      }
      for _, asset1 := range assets1 {
         fmt.Print("\n", &asset1, "\n")
      }
      time.Sleep(time.Second)
   }
}

func TestFields(t *testing.T) {
   for _, test1 := range tests {
      var fields1 Fields
      err := fields1.New(test1)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", fields1.algolia_convert_tracking())
      time.Sleep(time.Second)
   }
}

var tests = []struct {
   id  string
   url string
}{
   {
      url: "canalplus.cz/stream/film/argylle-tajny-agent",
   },
   {
      id:  "XT0kyelnPAOl3f-Bx7etkj_yX3nDHom_ymdCRK5A",
      url: "canalplus.cz/stream/series/fbi",
   },
   {
      url: "canalplus.cz/stream/series/silo",
   },
}
