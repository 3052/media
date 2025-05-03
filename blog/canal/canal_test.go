package canal

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := assets()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   file, err := os.Create("canal.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
}
