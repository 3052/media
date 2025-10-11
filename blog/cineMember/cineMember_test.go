package cineMember

import (
   "os"
   "os/exec"
   "testing"
)

func TestWrite(t *testing.T) {
   user, err := output("credential", "-h", "cinemember.nl", "-k", "user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h", "cinemember.nl")
   if err != nil {
      t.Fatal(err)
   }
   var sessionVar session
   err = sessionVar.New()
   if err != nil {
      t.Fatal(err)
   }
   err = sessionVar.login(user, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache + "/cineMember/session", []byte(sessionVar.String()), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestRead(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/cineMember/session")
   if err != nil {
      t.Fatal(err)
   }
   var sessionVar session
   err = sessionVar.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   resp, err := sessionVar.stream()
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
