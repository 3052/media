package binge

import (
   "os"
   "os/exec"
   "strings"
   "testing"
   "time"
)

func Test(t *testing.T) {
   data, err := exec.Command("password", "binge.com.au").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var token1 token
   err = token1.New(username, password)
   if err != nil {
      t.Fatal(err)
   }
   time.Sleep(time.Second)
   resp, err := token1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
