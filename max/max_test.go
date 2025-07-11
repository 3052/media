package max

import (
   "fmt"
   "os"
   "testing"
)

var content_tests = []struct {
   url      string
   location []string
}{
   {
      location: []string{"united states"},
      url:      "max.com/movies/dune/e7dc7b3a-a494-4ef1-8107-f4308aa6bbf7",
   },
   {
      url: "max.com/movies/despicable-me-4/3b1e1236-d69f-49f8-88df-2f57ab3c3ac7",
      location: []string{
         "chile",
         "colombia",
         "indonesia",
         "malaysia",
         "mexico",
         "peru",
         "philippines",
         "singapore",
         "thailand",
      },
   },
   {
      url: "max.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505",
      location: []string{
         "belgium",
         "brazil",
         "bulgaria",
         "chile",
         "colombia",
         "croatia",
         "czech republic",
         "denmark",
         "finland",
         "france",
         "hungary",
         "indonesia",
         "malaysia",
         "mexico",
         "netherlands",
         "norway",
         "peru",
         "philippines",
         "poland",
         "portugal",
         "romania",
         "singapore",
         "slovakia",
         "spain",
         "sweden",
         "thailand",
         "united states",
      },
   },
}

func TestContent(t *testing.T) {
   for _, test := range content_tests {
      fmt.Println(test)
   }
}

func TestPlayReady(t *testing.T) {
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
   // max.com/movies/dune/e7dc7b3a-a494-4ef1-8107-f4308aa6bbf7
   data, err = login1.PlayReady("06a38397-862d-4419-be84-0641939825e7")
   if err != nil {
      t.Fatal(err)
   }
   var play Playback
   err = play.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      home + "/media/max/PlayReady",
      []byte(play.Drm.Schemes.PlayReady.LicenseUrl), os.ModePerm,
   )
}
