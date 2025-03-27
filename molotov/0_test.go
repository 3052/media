package molotov

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestLogin(t *testing.T) {
   data, err := exec.Command("password", "molotov.tv").Output()
   if err != nil {
      panic(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   data, err = new_login(email, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/molotov/login", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
