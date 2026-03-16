package tubi

import "testing"

var tests = []struct {
   url string
   drm bool
}{
   {
      url: "https://tubitv.com/movies/617502",
      drm: true,
   },
   {
      url: "https://tubitv.com/series/300015509",
      drm: true,
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
