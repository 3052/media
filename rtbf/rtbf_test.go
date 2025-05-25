package rtbf

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   path string
   url  string
}{
   {
      path: "/emission/thelma-et-louise-29388",
      url:  "auvio.rtbf.be/emission/thelma-et-louise-29388",
   },
   {
      path: "/media/agatha-christie-pourquoi-pas-evans-agatha-christie-pourquoi-pas-evans-3280380",
      url:  "auvio.rtbf.be/media/agatha-christie-pourquoi-pas-evans-agatha-christie-pourquoi-pas-evans-3280380",
   },
}

func TestPage(t *testing.T) {
   for _, test := range tests {
      content1, err := Address{test.path}.Content()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", content1)
      time.Sleep(time.Second)
   }
}
