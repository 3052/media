package hulu

import "testing"

func Test(t *testing.T) {
   t.Log(tests)
}

var tests = []struct {
   url     string
   quality string
}{
   {
      url:     "https://hulu.com/movie/stay-5742941d-4b4a-4914-8774-f5d8d57f9382",
      quality: "2160p",
   },
   {
      url:     "https://hulu.com/movie/palm-springs-f70dfd4d-dbfb-46b8-abb3-136c841bba11",
      quality: "1080p",
   },
   {
      url: "https://hulu.com/series/house-ef39603f-eb90-4248-8237-f6168d7c1be1",
   },
}
