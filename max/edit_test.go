package max

import (
   "fmt"
   "log"
   "net/http"
   "os"
   "testing"
   "time"
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
   var line bool
   for season1, err := range login1.seasons() {
      if err != nil {
         t.Fatal(err)
      }
      for _, episode := range season1.sorted() {
         if line {
            fmt.Println()
         } else {
            line = true
         }
         fmt.Println(&episode)
      }
      time.Sleep(99*time.Millisecond)
   }
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}
func TestMovie(t *testing.T) {
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
   movie, err := login1.movie("/movie/12199308-9afb-460b-9d79-9d54b5d2514c")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(movie.movie())
}
