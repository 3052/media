package binge

import (
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/Auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 Auth
   err = auth1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(auth1.AccessToken, "\n\n")
   err = auth1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(auth1.AccessToken)
}

func TestWrite(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := exec.Command("password", "binge.com.au").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   data, err = NewAuth(username, password)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/binge/Auth", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
