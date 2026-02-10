package amc

import (
   "fmt"
   "testing"
)

func TestLicense(t *testing.T) {
   fmt.Println(key_tests)
}

var key_tests = []struct{
   key_id string
   url string
}{
   {
      key_id: "+7nUc5piRu2GY3lAiA4MvQ==",
      url: "amcplus.com/movies/nocebo--1061554",
   },
   {
      key_id: "vHkdO0RPSsqD3iPzeupPeA==",
      url: "amcplus.com/shows/orphan-black/episodes/season-1-instinct--1011152",
   },
}

var path_tests = []string{
   "/movies/nocebo--1061554",
   "amcplus.com/movies/nocebo--1061554",
   "https://www.amcplus.com/movies/nocebo--1061554",
   "www.amcplus.com/movies/nocebo--1061554",
}

func TestPath(t *testing.T) {
   for _, test := range path_tests {
      var web Address
      err := web.Set(test)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(web)
   }
}
