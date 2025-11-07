package mubi

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   fmt.Println(tests)
}

var tests = []struct {
   id        int64
   url       string
   locations []string
}{
   {
      id:     325455,
      url:    "mubi.com/films/passages-2022",
      locations: []string{
         "Austria",
         "Belgium",
         "Brazil",
         "Canada",
         "Chile",
         "Colombia",
         "Germany",
         "Ireland",
         "Italy",
         "Mexico",
         "Netherlands",
         "Peru",
         "Turkey",
         "United Kingdom",
         "United States",
      },
   },
   {
      url: "mubi.com/films/perfect-days",
      locations: []string{
         "Austria",
         "Brazil",
         "Chile",
         "Colombia",
         "Germany",
         "Ireland",
         "Mexico",
         "Peru",
         "Turkey",
         "United Kingdom",
      },
   },
}
