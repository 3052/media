package cineMember

import (
   "fmt"
   "os"
   "os/exec"
   "testing"
)

var tests = []struct {
   play  string
   title string
}{
   {
      play:  "https://cinemember.nl/nl/title/906945/american-hustle/play",
      title: "https://cinemember.nl/nl/title/468545/american-hustle",
   },
   { // buffer too small
      play:  "https://cinemember.nl/nl/title/904309/knives-out/play",
      title: "https://cinemember.nl/nl/title/469991/knives-out",
   },
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
   idVar, err := id(tests[0].title)
   if err != nil {
      t.Fatal(err)
   }
   streamVar, err := sessionVar.stream(idVar)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(streamVar.mpd())
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
