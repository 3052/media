package canal

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestOne(t *testing.T) {
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
   resp, err := ticket1.one(username, password)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
