package amc

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var tests = []struct {
   address string
   key string
   pssh string
} {
   // amcplus.com/shows/orphan-black/episodes/season-1-instinct--1011152
   episode: {
      address: "/shows/orphan-black/episodes/season-1-instinct--1011152",
      key: "95f11e40064f47007e7d950bd52d7b95",
      pssh: "AAAAVnBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAADYIARIQJqlCz6NjSI2kDWew20wbGRoNd2lkZXZpbmVfdGVzdCIIMTIzNDU2NzgyB2RlZmF1bHQ=",
   },
   // amcplus.com/movies/nocebo--1061554
   movie: {address: "/movies/nocebo--1061554"},
}

func Test_Content(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   auth, err := Read_Auth(home + "/mech/amc.json")
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests {
      con, err := auth.Content(test.address)
      if err != nil {
         t.Fatal(err)
      }
      video, err := con.Video_Player()
      if err != nil {
         t.Fatal(err)
      }
      name, err := video.Name()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      time.Sleep(time.Second)
   }
}

const (
   episode = iota
   movie
)