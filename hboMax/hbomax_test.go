package hboMax

import (
   "encoding/xml"
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   dir, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(dir + "/hboMax/userCache.xml")
   if err != nil {
      t.Fatal(err)
   }
   var cache struct {
      Login Login
   }
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      t.Fatal(err)
   }
   // hbomax.com/movies/dune/e7dc7b3a-a494-4ef1-8107-f4308aa6bbf7
   play, err := cache.Login.PlayReady("06a38397-862d-4419-be84-0641939825e7")
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      dir+"/hboMax/PlayReady",
      []byte(play.Drm.Schemes.PlayReady.LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

var content_tests = []struct {
   url      string
   location []string
}{
   {
      url:      "https://hbomax.com/movies/love-lies-bleeding/552f0116-65dc-4a87-9666-b2ef6135ed3d",
      location: []string{"united states"},
   },
   {
      url: "https://hbomax.com/movies/despicable-me-4/3b1e1236-d69f-49f8-88df-2f57ab3c3ac7",
      location: []string{
         "chile",
         "colombia",
         "indonesia",
         "malaysia",
         "mexico",
         "peru",
         "philippines",
         "singapore",
         "thailand",
      },
   },
   {
      url: "https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505",
      location: []string{
         "belgium",
         "brazil",
         "bulgaria",
         "chile",
         "colombia",
         "croatia",
         "czech republic",
         "denmark",
         "finland",
         "france",
         "hungary",
         "indonesia",
         "malaysia",
         "mexico",
         "netherlands",
         "norway",
         "peru",
         "philippines",
         "poland",
         "portugal",
         "romania",
         "singapore",
         "slovakia",
         "spain",
         "sweden",
         "thailand",
         "united states",
      },
   },
}

func TestContent(t *testing.T) {
   t.Log(content_tests)
}
