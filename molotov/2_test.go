package molotov

import (
   "os"
   "testing"
)

func TestView(t *testing.T) {
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
   resp, err := token.view()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
