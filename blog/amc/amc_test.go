package amc

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := series_detail()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   file, err := os.Create("amc.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
}
