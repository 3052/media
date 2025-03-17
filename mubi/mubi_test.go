package mubi

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var tests = []struct {
   id        int64
   key_id    string
   url       string
   locations []string
}{
   {
      id:     325455,
      key_id: "CA215A25BB2D43F0BD095FC671C984EE",
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

func Test(t *testing.T) {
   data, err := os.ReadFile("authenticate.txt")
   if err != nil {
      t.Fatal(err)
   }
   var auth Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, test1 := range tests {
      data, err := auth.SecureUrl(&Film{Id: test1.id})
      if err != nil {
         t.Fatal(err)
      }
      var secure SecureUrl
      err = secure.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", secure)
      time.Sleep(time.Second)
   }
}
