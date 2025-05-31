package rakuten

import (
   "fmt"
   "os"
   "testing"
)

var test = struct {
   season string
   url    string
}{
   season: "clink-1",
   url:    "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
}

func TestSeason(t *testing.T) {
   var web address
   err := web.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("address", []byte(test.url), os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
   seasons, err := web.seasons()
   if err != nil {
      t.Fatal(err)
   }
   for _, season1 := range seasons {
      fmt.Println(&season1)
   }
}
