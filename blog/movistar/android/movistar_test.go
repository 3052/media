package movistar

import (
   "fmt"
   "testing"
)

var test = struct {
   id  int64
   url string
}{
   id:  3427440,
   url: "movistarplus.es/cine/ficha?id=3427440",
}

func Test(t *testing.T) {
   var details1 details
   err := details1.New(test.id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", details1)
}
