package sky

import (
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestWrite(t *testing.T) {
   data, err := exec.Command("password", "sky.ch").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var login1 login
   err = login1.New()
   if err != nil {
      t.Fatal(err)
   }
   cookies1, err := login1.login(username, password)
   if err != nil {
      t.Fatal(err)
   }
   session, ok := cookies1.session_id()
   if !ok {
      t.Fatal(".session_id()")
   }
   file, err := os.Create("session_id")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   fmt.Fprint(file, session)
}

func TestRead(t *testing.T) {
   data, err := os.ReadFile("session_id")
   if err != nil {
      t.Fatal(err)
   }
   var session cookie
   err = session.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(session)
}

func TestWeb(t *testing.T) {
   data, err := os.ReadFile("session_id")
   if err != nil {
      t.Fatal(err)
   }
   var session cookie
   err = session.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   data, err = sky_player(session[0])
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(data))
}
