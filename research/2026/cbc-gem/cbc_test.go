package cbc

import (
   "fmt"
   "os"
   "testing"
)

func TestProfile(t *testing.T) {
   var token LoginToken
   username, password := os.Getenv("cbc_username"), os.Getenv("cbc_password")
   if err := token.New(username, password); err != nil {
      t.Fatal(err)
   }
   profile, err := token.Profile()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(profile)
}
