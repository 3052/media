package paramount

import "testing"

var resolved = []string{
   "https://paramountplus.com/ie/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ/",
   "https://paramountplus.com/ie/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ/",
   "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q/",
   "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD/",
}

func TestPath(t *testing.T) {
   path, err := FetchPath(
      "http://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   )
   if err != nil {
      t.Fatal(err)
   }
   t.Log(path)
}

func TestResolve(t *testing.T) {
   t.Log(resolved)
}
