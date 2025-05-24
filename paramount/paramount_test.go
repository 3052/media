package paramount

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   content_id string
   location   []string
   url        string
}{
   {
      content_id: "tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      url:        "paramountplus.com/movies/video/tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      location:   []string{"USA"},
   },
   {
      content_id: "esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      location:   []string{"USA"},
      url:        "paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
   {
      content_id: "rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      location:   []string{"USA"},
      url:        "cbs.com/shows/video/rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
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
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      at1, err := ComCbsApp.At()
      if err != nil {
         t.Fatal(err)
      }
      session1, err := at1.Session(test1.content_id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(session1)
      time.Sleep(time.Second)
   }
}
