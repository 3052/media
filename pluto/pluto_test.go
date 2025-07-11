package pluto

import (
   "fmt"
   "testing"
)

var tests = []struct {
   id  string
   url string
}{
   {
      id:  "5a9dd73dfb6f2f17481aff11",
      url: "pluto.tv/us/on-demand/movies/5a9dd73dfb6f2f17481aff11",
   },
   {
      id:  "66d0bb64a1c89200137fb0e6",
      url: "pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      video, err := NewVod(testVar.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", video)
   }
}
