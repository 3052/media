package paramount

import "testing"

var resolved = []string{
   // pass
   "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q/",
   "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD/",
   // fail
   "https://paramountplus.com/ie/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ/",
   "https://paramountplus.com/ie/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ/",
}

func TestResolve(t *testing.T) {
   t.Log(resolved)
}
