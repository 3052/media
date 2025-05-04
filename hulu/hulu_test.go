package hulu

import (
   "os"
   "path"
   "testing"
   "time"
)

var tests = []struct {
   content string
   url     string
}{
   {
      content: "episode",
      url:     "hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d",
   },
   {
      content: "film",
      url:     "hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11",
   },
}

func TestDeepLink(t *testing.T) {
   for _, test1 := range tests {
      data, err := os.ReadFile("authenticate.txt")
      if err != nil {
         t.Fatal(err)
      }
      var auth Authenticate
      err = auth.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      base := path.Base(test1.url)
      _, err = auth.DeepLink(Entity{base})
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
