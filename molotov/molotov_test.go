package molotov

import (
   "os"
   "testing"
)

func TestWidevine(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/molotov/refresh")
   if err != nil {
      t.Fatal(err)
   }
   var token Refresh
   err = token.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   var web Address
   err = web.Set("molotov.tv/fr_fr/p/15082-531")
   if err != nil {
      t.Fatal(err)
   }
   _, err = token.View(&web)
   if err != nil {
      t.Fatal(err)
   }
}
