package movistar

import (
   "fmt"
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
   var token1 token
   err = token1.New(username, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(token1.AccessToken, "\n", token1.duration(), "\n")
}
