package movistar

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestToken(t *testing.T) {
   data, err := exec.Command("password", "movistarplus.es").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   data, err = new_token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/movistar/token", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
