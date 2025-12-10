package pluto

import (
   "io"
   "net/http"
   "strings"
   "testing"
)

func Test(t *testing.T) {
   var value Series
   err := value.Fetch("6495eff09263a40013cf63a5")
   if err != nil {
      t.Fatal(err)
   }
   resp, err := http.Get(value.String())
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   const height = ` height="1080"`
   if !strings.Contains(string(data), height) {
      t.Fatal(height)
   }
}
