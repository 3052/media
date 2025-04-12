package movistar

import (
   "os"
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
   resp, err := details(test.id)
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   file, err := os.Create(".json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   file.ReadFrom(resp.Body)
}
