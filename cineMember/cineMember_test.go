package cineMember

import (
   "fmt"
   "os"
   "os/exec"
   "testing"
)

var tests = []string{
   "https://cinemember.nl/nl/title/468845/the-worst-person-in-the-world",
   "https://cinemember.nl/nl/title/469991/knives-out", // buffer too small
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
   var sessionVar Session
   err = sessionVar.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   idVar, err := Id(tests[0])
   if err != nil {
      t.Fatal(err)
   }
   streamVar, err := sessionVar.Stream(idVar)
   if err != nil {
      t.Fatal(err)
   }
   for _, link := range streamVar.Links {
      fmt.Printf("%+v\n", link)
   }
}

func TestWrite(t *testing.T) {
   user, err := output("credential", "-h", "cinemember.nl", "-k", "user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h", "cinemember.nl")
   if err != nil {
      t.Fatal(err)
   }
   var sessionVar Session
   err = sessionVar.New()
   if err != nil {
      t.Fatal(err)
   }
   err = sessionVar.Login(user, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/cineMember/session", []byte(sessionVar.String()), os.ModePerm,
   )
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
