package disney

import (
   "fmt"
   "os/exec"
   "testing"
)

func TestLogin(t *testing.T) {
   email, err := output("credential", "-h=disneyplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=disneyplus.com")
   if err != nil {
      t.Fatal(err)
   }
   token, err := register_device()
   if err != nil {
      t.Fatal(err)
   }
   account, err := token.login(email, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", account)
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}
