package tubi

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   drm      bool
   location string
   url      string
}{
   {
      url: "https://tubitv.com/movies/100047876",
      drm: true,
   },
   {
      url: "tubitv.com/tv-shows/200042567",
      drm: true,
   },
   {
      url: "tubitv.com/movies/667315",
      drm: false,
   },
   {
      location: "Australia",
      url:      "tubitv.com/movies/643397",
      drm:      false,
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
