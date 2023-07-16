package amc

import (
   "154.pages.dev/media"
   "fmt"
   "os"
   "testing"
   "time"
)

func Test_Content(t *testing.T) {
   var auth Auth_ID
   {
      s, err := os.UserHomeDir()
      if err != nil {
         t.Fatal(err)
      }
      b, err := os.ReadFile(s + "/amc/auth.json")
      if err != nil {
         t.Fatal(err)
      }
      auth.Unmarshal(b)
   }
   for _, test := range tests {
      con, err := auth.Content(test.address)
      if err != nil {
         t.Fatal(err)
      }
      vid, err := con.Video()
      if err != nil {
         t.Fatal(err)
      }
      name, err := media.Name(vid)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      time.Sleep(time.Second)
   }
}

func Test_Refresh(t *testing.T) {
   var auth Auth_ID
   {
      s, err := os.UserHomeDir()
      if err != nil {
         t.Fatal(err)
      }
      b, err := os.ReadFile(s + "/amc/auth.json")
      if err != nil {
         t.Fatal(err)
      }
      auth.Unmarshal(b)
   }
   if err := auth.Refresh(); err != nil {
      t.Fatal(err)
   }
   {
      b, err := auth.Marshal()
      if err != nil {
         t.Fatal(err)
      }
      os.WriteFile(os.Getenv("AMC_PLUS"), b, 0666)
   }
}

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

const (
   episode = iota
   movie
)
