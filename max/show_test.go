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
   items, err := login1.show("/show/14f9834d-bc23-41a8-ab61-5c8abdbea505")
   if err != nil {
      t.Fatal(err)
   }
   season, ok := items.season()
   fmt.Printf("%+v %v\n", season, ok)
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Print(req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}
