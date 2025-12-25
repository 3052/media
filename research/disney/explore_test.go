package disney

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestExplore(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/account.xml")
   if err != nil {
      t.Fatal(err)
   }
   var account_with account
   err = xml.Unmarshal(data, &account_with)
   if err != nil {
      t.Fatal(err)
   }
   explore, err := account_with.explore(test.entity)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(explore.restart())
}

var test = struct {
   entity string
   url    string
}{
   entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   url:    "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
}
