package canal

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   id  string
   url string
}{
   {
      id:  "XT0kyelnPAOl3f-Bx7etkj_yX3nDHom_ymdCRK5A",
      url: "canalplus.cz/stream/series/fbi",
   },
   {
      id:  "ZXkaWHVpx827Fz_4ZNtW5l8MoKD5_2lhv0nYe4m3",
      url: "canalplus.cz/stream/series/mozart-v-dzungli",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      fmt.Println(test1.url)
      assets1, err := assets(test1.id, 1)
      if err != nil {
         t.Fatal(err)
      }
      for _, asset1 := range assets1 {
         fmt.Printf("%+v\n", asset1)
      }
      time.Sleep(time.Second)
   }
}
