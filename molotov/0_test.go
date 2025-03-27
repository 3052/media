package molotov

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func Test(t *testing.T) {
   data, err := exec.Command("password", "molotov.tv").Output()
   if err != nil {
      panic(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   resp, err := zero(email, password)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
