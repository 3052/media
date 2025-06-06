package max

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   for _, test1 := range tests {
      fmt.Println(test1)
   }
}

var tests = []struct {
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
