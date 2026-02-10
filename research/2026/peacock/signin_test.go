package peacock

import (
   "fmt"
   "log/slog"
   "os"
   "testing"
)

func TestSignWrite(t *testing.T) {
   slog.SetLogLoggerLevel(slog.LevelDebug)
   user, password := os.Getenv("peacock_username"), os.Getenv("peacock_password")
   if user == "" {
      t.Fatal("peacock_username")
   }
   var sign SignIn
   err := sign.New(user, password)
   if err != nil {
      t.Fatal(err)
   }
   text, err := sign.Marshal()
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile(home + "/peacock.json", text, 0666)
}

func TestSignRead(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile(home + "/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(text)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", auth)
}
