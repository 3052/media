package paramount

import "testing"

var videos = []struct {
   resolution string
   url        string
}{
   {
      resolution: "2160p",
      url:        "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      resolution: "2160p",
      url:        "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
   {
      resolution: "2160p",
      url:        "https://paramountplus.com/gb/movies/video/yh8qG9949D8fvJa0dFmY1C_SilOYt2hS",
   },
}

func Test(t *testing.T) {
   t.Log(videos)
}
