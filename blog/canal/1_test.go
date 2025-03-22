package canal

import (
   "fmt"
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
   token1, err := ticket1.token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\nn", token1)
}
