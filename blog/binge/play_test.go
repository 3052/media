package binge

import (
   "os"
   "testing"
)

func TestPlay(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 auth
   err = auth1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token, err := auth1.token()
   if err != nil {
      t.Fatal(err)
   }
   resp, err := token.play()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
