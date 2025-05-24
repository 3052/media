package hulu

import (
   "fmt"
   "testing"
)

var tests = []struct {
   content string
   url     string
}{
   {
      content: "film",
      url:     "hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11",
   },
   {
      content: "episode",
      url:     "hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d",
   },
}

func TestDeepLink(t *testing.T) {
   fmt.Println(tests)
}
