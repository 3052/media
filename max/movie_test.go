package max

import (
   "fmt"
   "os"
   "testing"
)

func TestItems(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/max/Login")
   if err != nil {
      t.Fatal(err)
   }
   var login1 Login
   err = login1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   items1, err := login1.items("/movie/12199308-9afb-460b-9d79-9d54b5d2514c")
   if err != nil {
      t.Fatal(err)
   }
   for episode := range items1.episode() {
      fmt.Printf("%+v\n\n", episode)
   }
}
