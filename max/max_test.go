package max

import (
   "fmt"
   "os"
   "testing"
   "time"
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
   for movie := range items.Seq() {
      fmt.Println(&movie)
   }
}

func TestPlayback(t *testing.T) {
   data, err := os.ReadFile("login.txt")
   if err != nil {
      t.Fatal(err)
   }
   var login1 Login
   err = login1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, test1 := range tests {
      data, err := login1.Playback(test1.one)
      if err != nil {
         t.Fatal(err)
      }
      var play Playback
      err = play.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", play.Fallback)
      time.Sleep(time.Second)
   }
}

var tests = []struct {
   zero     string
   one      string
   location []string
}{
   {
      zero: "play.max.com/movie/83e518fa-7f76-47d0-a607-227b53bf3e6c",
      location: []string{
         "Brazil",
         "Chile",
         "Colombia",
         "Denmark",
         "Finland",
         "France",
         "Mexico",
         "Norway",
         "Peru",
         "Sweden",
         "United States",
      },
   },
   {
      zero: "play.max.com/show/14f9834d-bc23-41a8-ab61-5c8abdbea505",
      location: []string{
         "Belgium",
         "Brazil",
         "Bulgaria",
         "Chile",
         "Colombia",
         "Croatia",
         "Czech Republic",
         "Denmark",
         "Finland",
         "France",
         "Hungary",
         "Indonesia",
         "Malaysia",
         "Mexico",
         "Netherlands",
         "Norway",
         "Peru",
         "Philippines",
         "Poland",
         "Portugal",
         "Romania",
         "Singapore",
         "Slovakia",
         "Spain",
         "Sweden",
         "Thailand",
         "United States",
      },
   },
   {
      zero: "play.max.com/movie/3b1e1236-d69f-49f8-88df-2f57ab3c3ac7",
      location: []string{
         "Chile",
         "Colombia",
         "Indonesia",
         "Malaysia",
         "Mexico",
         "Peru",
         "Philippines",
         "Singapore",
         "Thailand",
      },
   },
}
