package gem

import (
   "fmt"
   "testing"
)

func Test_New_Profile(t *testing.T) {
   u, err := user(home + "/gem.cbc.ca/user.json")
   if err != nil {
      t.Fatal(err)
   }
   tok, err := New_Token(u["username"], u["password"])
   if err != nil {
      t.Fatal(err)
   }
   pro, err := tok.Profile()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", pro)
}

func Test_Profile(t *testing.T) {
   u, err := user(home + "/gem.cbc.ca/user.json")
   if err != nil {
      t.Fatal(err)
   }
   tok, err := New_Token(u[0], u[1])
   if err != nil {
      t.Fatal(err)
   }
   pro, err := tok.Profile()
   if err != nil {
      t.Fatal(err)
   }
   if err := pro.Write_File(home + "/cbc.json"); err != nil {
      t.Fatal(err)
   }
}