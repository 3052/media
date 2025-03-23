package max

import (
   "fmt"
   "os"
   "testing"
)

func TestRoutes(t *testing.T) {
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
   routes1, err := login1.routes("/movie/12199308-9afb-460b-9d79-9d54b5d2514c")
   if err != nil {
      t.Fatal(err)
   }
   edit1, ok := routes1.edit()
   fmt.Printf("%+v %v\n", edit1, ok)
}
