package canal

import (
   "fmt"
   "os"
   "testing"
)

func TestToken(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/canal/sso_token")
   if err != nil {
      t.Fatal(err)
   }
   var sso sso_token
   err = sso.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token1, err := sso.token()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", token1)
}
