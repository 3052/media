package max

import (
   "fmt"
   "log"
   "net/http"
   "os"
   "testing"
   "time"
)

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
   movie, err := login1.movie("12199308-9afb-460b-9d79-9d54b5d2514c")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(movie.movie())
}

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
   season, err := login1.season("14f9834d-bc23-41a8-ab61-5c8abdbea505", 1)
   if err != nil {
      t.Fatal(err)
   }
   season1 := slices.SortedFunc(season.episode(), func(a, b video) int {
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   for i, episode := range season1 {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
}
