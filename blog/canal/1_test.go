package canal

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestToken(t *testing.T) {
   data, err := exec.Command("password", "canalplus.cz").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var ticket1 ticket
   err = ticket1.New()
   if err != nil {
      t.Fatal(err)
   }
   data, err = ticket1.token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/canal/token", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
