package mubi

import (
   "os"
   "testing"
)

func TestAuthenticate(t *testing.T) {
   var (
      code LinkCode
      err error
   )
   code.Raw, err = os.ReadFile("code.txt")
   if err != nil {
      t.Fatal(err)
   }
   code.Unmarshal()
   auth, err := code.Authenticate()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("authenticate.txt", auth.Raw, os.ModePerm)
}
