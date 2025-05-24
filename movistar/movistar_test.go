package movistar

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

var test = struct {
   id  int64
   url string
}{
   id:  3427440,
   url: "movistarplus.es/cine/ficha?id=3427440",
}

func TestDetails(t *testing.T) {
   _, err := NewDetails(test.id)
   if err != nil {
      t.Fatal(err)
   }
}

func TestDevice(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/movistar/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 Token
   err = token1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   oferta1, err := token1.Oferta()
   if err != nil {
      t.Fatal(err)
   }
   data, err = token1.Device(oferta1)
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
   data, err = NewToken(username, password)
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
