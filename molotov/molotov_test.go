package molotov

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestRefresh(t *testing.T) {
   data, err := exec.Command("password", "molotov.tv").Output()
   if err != nil {
      t.Fatal(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   var login1 login
   err = login1.New(email, password)
   if err != nil {
      t.Fatal(err)
   }
   data, err = login1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/molotov/refresh", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
