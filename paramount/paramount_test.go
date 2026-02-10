package paramount

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   content_id string
   key_id     string
   location   []string
   url        string
}{
   {
      content_id: "WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      key_id:     "bsT01+Q1Ta+39TayayKhBg==",
      url:        "paramountplus.com/shows/video/WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      location:   []string{"Australia"},
   },
   {
      content_id: "Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      key_id:     "BsO37qHORXefruKryNAaVQ==",
      url:        "paramountplus.com/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      location:   []string{"Australia", "United Kingdom"},
   },
   {
      content_id: "3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      url: "paramountplus.com/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      location: []string{
         "Brazil", "Canada", "Chile", "Colombia", "Mexico", "Peru",
      },
   },
   {
      content_id: "tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      url:        "paramountplus.com/movies/video/tOeI0WHG3icuPhCk5nkLXNmi5c4Jfx41",
      location:   []string{"USA"},
   },
   {
      content_id: "esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      key_id:     "H94BVNcqT0WRKzTwzgd36w==",
      location:   []string{"USA"},
      url:        "paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
   {
      content_id: "rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      key_id:     "Sryog4HeT2CLHx38NftIMA==",
      location:   []string{"USA"},
      url:        "cbs.com/shows/video/rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      session, err := ComCbsApp.Session(test1.content_id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(session)
      time.Sleep(time.Second)
   }
}
