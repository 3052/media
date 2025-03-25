package max

import (
   "fmt"
   "log"
   "net/http"
   "os"
   "slices"
   "testing"
)

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
   items, err := login1.Movie("12199308-9afb-460b-9d79-9d54b5d2514c")
   if err != nil {
      t.Fatal(err)
   }
   for movie := range items.Movie() {
      fmt.Println(&movie)
   }
}

func TestSeason(t *testing.T) {
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
   items, err := login1.Season("14f9834d-bc23-41a8-ab61-5c8abdbea505", 1)
   if err != nil {
      t.Fatal(err)
   }
   episodes := slices.SortedFunc(items.Episode(), func(a, b Video) int {
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}
