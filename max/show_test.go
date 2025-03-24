package max

import (
   "fmt"
   "log"
   "net/http"
   "os"
   "testing"
)

func TestShow(t *testing.T) {
   http.DefaultClient.Transport = transport{}
   log.SetFlags(log.Ltime)
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/max/Login")
   if err != nil {
      t.Fatal(err)
   }
   var login1 Login
   err = login1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   season1, err := login1.season()
   if err != nil {
      t.Fatal(err)
   }
   for i, episode := range season1.episode() {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Print(req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}
