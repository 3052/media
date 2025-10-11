package cineMember

import (
   "os"
   "os/exec"
   "testing"
)

func Test(t *testing.T) {
   user, err := output("credential", "-h", "cinemember.nl", "-k", "user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h", "cinemember.nl")
   if err != nil {
      t.Fatal(err)
   }
   cookie, err := session()
   if err != nil {
      t.Fatal(err)
   }
   err = login(cookie, user, password)
   if err != nil {
      t.Fatal(err)
   }
   resp, err := stream(cookie)
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
