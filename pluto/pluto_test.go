package pluto

import (
   "fmt"
   "testing"
)

var tests = []struct{
   id     string
   key_id string
   url    string
}{
   {
      id:     "5c4bb2b308d10f9a25bbc6af",
      key_id: "AAAAAGbZBRrrxvnmpuNLhg==",
      url:    "pluto.tv/on-demand/movies/5c4bb2b308d10f9a25bbc6af",
   },
   {
      url:    "pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6/episode/66fb16fda2922a00135e87f7",
      id:     "66fb16fda2922a00135e87f7",
      key_id: "",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      var web Address
      err := web.Set(test1.url)
      if err != nil {
         t.Fatal(err)
      }
      video, err := web.Vod("")
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", video)
   }
}
