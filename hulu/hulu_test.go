package hulu

import (
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   //data, err := os.ReadFile(cache + "/hulu/Session")
   if err != nil {
      t.Fatal(err)
   }
   var session_var Session
   //err = session_var.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   err = session_var.TokenRefresh()
   if err != nil {
      t.Fatal(err)
   }
   test := tests[0]
   play, err := session_var.Playlist(&DeepLink{EabId: test.id})
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/hulu/PlayReady", []byte(play.DashPrServer), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

var tests = []struct {
   id      string
   url     string
   quality string
}{
   {
      id:      "EAB::f70dfd4d-dbfb-46b8-abb3-136c841bba11::61556664::101167038",
      url:     "https://hulu.com/movie/palm-springs-f70dfd4d-dbfb-46b8-abb3-136c841bba11",
      quality: "1080p",
   },
   {
      url:     "https://hulu.com/movie/stay-5742941d-4b4a-4914-8774-f5d8d57f9382",
      quality: "2160p",
   },
   {
      url: "https://hulu.com/series/house-ef39603f-eb90-4248-8237-f6168d7c1be1",
   },
}
