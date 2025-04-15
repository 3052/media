package movistar

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestDevice(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/movistar/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   oferta1, err := token1.oferta()
   if err != nil {
      t.Fatal(err)
   }
   data, err = token1.device(oferta1)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/movistar/device", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}

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
