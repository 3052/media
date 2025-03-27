package molotov

import (
   "fmt"
   "os"
   "testing"
)

func TestAssets(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/molotov/refresh")
   if err != nil {
      t.Fatal(err)
   }
   var token refresh
   err = token.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   view1, err := token.view()
   if err != nil {
      t.Fatal(err)
   }
   assets1, err := token.assets(view1)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", assets1)
}
