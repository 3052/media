package movistar

import (
   "fmt"
   "os"
   "testing"
)

func TestDevices(t *testing.T) {
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
   device1, err := token1.device()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", device1)
}
