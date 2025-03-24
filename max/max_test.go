package max

import (
   "fmt"
   "os"
   "testing"
   "time"
)

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
      var watch WatchUrl
      err = watch.Set(test1.one)
      if err != nil {
         t.Fatal(err)
      }
      data, err := login1.Playback(&watch)
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
      one:  "play.max.com/video/watch/5c762883-279e-40ed-ab84-43fdda9d88a0",
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
      one:  "play.max.com/video/watch/28ae9450-8192-4277-b661-e76eaad9b2e6",
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
      one:  "play.max.com/video/watch/857fc45b-5652-42ca-9192-ac1e5e456300",
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
