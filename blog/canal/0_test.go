package canal

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   var login1 login
   err := login1.New()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", login1)
}
