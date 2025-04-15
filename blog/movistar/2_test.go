package movistar

import (
   "os"
   "testing"
)

func TestDevice(t *testing.T) {
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
   data, err = token1.device(oferta1)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/movistar/device", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
