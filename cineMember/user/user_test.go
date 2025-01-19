package user

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestUser(t *testing.T) {
   data, err := exec.Command("password", "cinemember.nl").Output()
   if err != nil {
      t.Fatal(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   data, err = marshal(email, password)
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("user.txt", data, os.ModePerm)
}