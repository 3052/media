package peacock

import (
   "fmt"
   "log/slog"
   "os"
   "testing"
)

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestSignWrite(t *testing.T) {
   user, err := output("credential", "-h=peacocktv.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=peacocktv.com")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   err := sign.New(user, password)
   if err != nil {
      t.Fatal(err)
   }
   data, err := sign.Marshal()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("peacock.json", data, 0666)
}

func TestSignRead(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(data)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", auth)
}
