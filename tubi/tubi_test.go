package tubi

import "testing"

var tests = []struct {
   url string
   drm bool
}{
   {
      url: "https://tubitv.com/movies/667315",
      drm: false,
   },
   {
      url: "https://tubitv.com/movies/100047876",
      drm: true,
   },
   {
      url: "https://tubitv.com/tv-shows/200042567",
      drm: true,
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
