package plex

import (
   "fmt"
   "testing"
   "time"
)

var watch_tests = []struct {
   drm  bool
   path string
   url  string
}{
   {
      drm:  true,
      path: "/movie/ghost-in-the-shell",
      url:  "watch.plex.tv/watch/movie/ghost-in-the-shell",
   },
   {
      url:  "watch.plex.tv/movie/limitless",
      path: "/movie/limitless",
   },
   {
      path: "/show/broadchurch/season/3/episode/5",
      url:  "watch.plex.tv/show/broadchurch/season/3/episode/5",
   },
}

func Test(t *testing.T) {
   data, err := NewUser()
   if err != nil {
      t.Fatal(err)
   }
   var user_var User
   err = user_var.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, watch_test := range watch_tests {
      match_var, err := user_var.Match(watch_test.path)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(match_var)
      time.Sleep(time.Second)
   }
}
