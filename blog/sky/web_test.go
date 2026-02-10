package sky

import (
   "fmt"
   "os"
   "testing"
)

func TestWeb(t *testing.T) {
   data, err := os.ReadFile("session.txt")
   if err != nil {
      t.Fatal(err)
   }
   var session Cookie
   err = session.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   data, err = sky_player(session.Cookie)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(data))
}
