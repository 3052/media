package amc

import (
   "154.pages.dev/text"
   "fmt"
   "os"
   "testing"
   "time"
)

func TestLogin(t *testing.T) {
   username := os.Getenv("amc_username")
   if username == "" {
      t.Fatal("Getenv")
   }
   password := os.Getenv("amc_password")
   var auth Authorization
   err := auth.Unauth()
   if err != nil {
      t.Fatal(err)
   }
   err = auth.Unmarshal(auth.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   err = auth.Login(username, password)
   if err != nil {
      t.Fatal(err)
   }
}

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   raw, err := os.ReadFile(home + "/amc.txt")
   if err != nil {
      t.Fatal(err)
   }
   var auth Authorization
   err = auth.Unmarshal(raw)
   if err != nil {
      t.Fatal(err)
   }
   err = auth.Refresh()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile(home + "/amc.txt", auth.Marshal(), 0666)
}

func TestContent(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   raw, err := os.ReadFile(home + "/amc.txt")
   if err != nil {
      t.Fatal(err)
   }
   var auth Authorization
   auth.Unmarshal(raw)
   for _, test := range tests {
      var web Address
      web.Set(test.url)
      content, err := auth.Content(web.Path)
      if err != nil {
         t.Fatal(err)
      }
      video, err := content.Video()
      if err != nil {
         t.Fatal(err)
      }
      name, err := text.Name(video)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", name)
      time.Sleep(time.Second)
   }
}

var path_tests = []string{
   "http://amcplus.com/movies/nocebo--1061554",
   "amcplus.com/movies/nocebo--1061554",
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
