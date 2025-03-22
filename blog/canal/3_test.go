package canal

import (
   "fmt"
   "os"
   "testing"
)

func TestPlay(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/canal/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   session1, err := token1.session()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := session1.play()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", play1)
}
