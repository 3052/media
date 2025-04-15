package movistar

import (
   "fmt"
   "log"
   "net/http"
   "os"
   "testing"
)

var test = struct {
   id  int64
   url string
}{
   id:  3427440,
   url: "movistarplus.es/cine/ficha?id=3427440",
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}

func TestSession(t *testing.T) {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = transport{}
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/movistar/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 token
   err = token1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   oferta1, err := token1.oferta()
   if err != nil {
      t.Fatal(err)
   }
   device1, err := token1.device(oferta1)
   if err != nil {
      t.Fatal(err)
   }
   init1, err := oferta1.init_data(device1)
   if err != nil {
      t.Fatal(err)
   }
   var details1 details
   err = details1.New(test.id)
   if err != nil {
      t.Fatal(err)
   }
   session1, err := device1.session(init1, &details1)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", session1)
}
