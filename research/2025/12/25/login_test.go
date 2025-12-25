package disney

import (
   "os"
   "os/exec"
   "testing"
)

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestLogin(t *testing.T) {
   email, err := output("credential", "-h=disneyplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=disneyplus.com")
   if err != nil {
      t.Fatal(err)
   }
   device, err := fetch_register_device()
   if err != nil {
      t.Fatal(err)
   }
   resp, err := device.login(email, password)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
