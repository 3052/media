package binge

import (
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestService(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(token1.AccessToken, "\n\n")
   err = token1.service()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(token1.AccessToken)
}

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(token1.AccessToken, "\n\n")
   err = token1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(token1.AccessToken)
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
   data, err = new_token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/binge/token", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
