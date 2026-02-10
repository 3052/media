package cbc

import (
   "154.pages.dev/rosso"
   "fmt"
   "os"
   "testing"
   "time"
)

var links = []string{
   "https://gem.cbc.ca/downton-abbey/s01e05",
   "https://gem.cbc.ca/the-fall/s02e03",
   "https://gem.cbc.ca/the-witch",
}

func TestStream(t *testing.T) {
   for _, link := range links {
      var gem GemCatalog
      err := gem.New(link)
      if err != nil {
         t.Fatal(err)
      }
      item, ok := gem.Item()
      if ok {
         fmt.Printf("%+v\n", item)
         fmt.Println(rosso.Name(gem.StructuredMetadata))
      }
      time.Sleep(time.Second)
   }
}

func TestMedia(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   var profile GemProfile
   profile.Raw, err = os.ReadFile(home + "/cbc/profile.json")
   if err != nil {
      t.Fatal(err)
   }
   for _, link := range links {
      var catalog GemCatalog
      err := catalog.New(link)
      if err != nil {
         t.Fatal(err)
      }
      item, ok := catalog.Item()
      if ok {
         media, err := profile.Media(item)
         if err != nil {
            t.Fatal(err)
         }
         fmt.Printf("%+v\n", media)
      }
      time.Sleep(time.Second)
   }
}
