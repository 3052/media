package molotov

import (
   "os"
   "testing"
)

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/molotov/login")
   if err != nil {
      t.Fatal(err)
   }
   var login1 login
   err = login1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   resp, err := login1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
