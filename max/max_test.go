package max

import (
   "fmt"
   "os"
   "testing"
)

func TestLogin(t *testing.T) {
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
   fmt.Printf("%+v\n", login1)
}

func TestContent(t *testing.T) {
   for _, test := range content_tests {
      fmt.Println(test)
   }
}

var content_tests = []struct {
   url      string
   location []string
}{
   {
      location: []string{"united states"},
      url:      "max.com/movies/dune/e7dc7b3a-a494-4ef1-8107-f4308aa6bbf7",
   },
   {
      url: "play.max.com/show/14f9834d-bc23-41a8-ab61-5c8abdbea505",
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
      url: "play.max.com/movie/3b1e1236-d69f-49f8-88df-2f57ab3c3ac7",
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
