package paramount

import (
   "fmt"
   "os"
   "testing"
)

var location_tests = []struct {
   content_id string
   location   []string
   url        string
}{
   {
      location:   []string{"USA"},
      url:        "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      content_id: "wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      content_id: "esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      location:   []string{"USA"},
      url:        "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
   {
      content_id: "Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      url:        "https://paramountplus.com/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      location:   []string{"Australia", "United Kingdom"},
   },
   {
      content_id: "3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      url:        "https://paramountplus.com/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      location: []string{
         "Brazil", "Canada", "Chile", "Colombia", "Mexico", "Peru",
      },
   },
}

func TestLocation(t *testing.T) {
   fmt.Println(location_tests)
}

func TestPlayReady(t *testing.T) {
   token, err := ComCbsApp.At()
   if err != nil {
      t.Fatal(err)
   }
   session_var, err := token.playReady("wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q")
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/paramount/PlayReady", []byte(session_var.LsSession), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}
