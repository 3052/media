package paramount

import "testing"

func TestLog(t *testing.T) {
   t.Log(contents)
}

var contents = [][]string{
   {
      "https://paramountplus.com/gb/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      "https://paramountplus.com/ie/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
   },
   {
      "https://paramountplus.com/gb/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      "https://paramountplus.com/ie/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
   },
   {
      "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      "https://paramountplus.com/movies/zodiac/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
}
