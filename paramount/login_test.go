package paramount

import (
   "os"
   "os/exec"
   "testing"
)

func TestLogin(t *testing.T) {
   user, err := output("credential", "-h=paramountplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=paramountplus.com")
   if err != nil {
      t.Fatal(err)
   }
   at, err := GetAt(ComCbsApp.AppSecret)
   if err != nil {
      t.Fatal(err)
   }
   resp, err := login(at, user, password)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}
