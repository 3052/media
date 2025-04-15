package movistar

import (
   "fmt"
   "os"
   "testing"
)

func TestSession(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/movistar/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   oferta1, err := token1.oferta()
   if err != nil {
      t.Fatal(err)
   }
   device1, err := token1.device(oferta1)
   if err != nil {
      t.Fatal(err)
   }
   init1, err := oferta1.init_data(device1)
   if err != nil {
      t.Fatal(err)
   }
   session1, err := device1.session(init1)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", session1)
}
