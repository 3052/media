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
   period int
}{
   {
      content_id: "rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      location:   []string{"USA"},
      url:        "cbs.com/shows/video/rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      period: 2,
   },
   {
      content_id: "3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      url:        "paramountplus.com/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      location: []string{
         "Brazil", "Canada", "Chile", "Colombia", "Mexico", "Peru",
      },
   },
   {
      content_id: "Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      url:        "paramountplus.com/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      location:   []string{"Australia", "United Kingdom"},
   },
   {
      content_id: "WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      url:        "paramountplus.com/shows/video/WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      location:   []string{"Australia"},
   },
   {
      content_id: "tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      url:        "paramountplus.com/movies/video/tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      location:   []string{"USA"},
      period: 1,
   },
   {
      content_id: "esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      location:   []string{"USA"},
      url:        "paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      period: 1,
   },
}

func TestPlayReady(t *testing.T) {
   token, err := ComCbsApp.At()
   if err != nil {
      t.Fatal(err)
   }
   sess, err := token.playReady("tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41")
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      home + "/media/paramount/PlayReady",
      []byte(sess.LsSession), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestLocation(t *testing.T) {
   fmt.Println(location_tests)
}
