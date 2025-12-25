package disney

import (
   "encoding/xml"
   "io"
   "os"
   "testing"
)

func TestExplore(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/refresh_token.xml")
   if err != nil {
      t.Fatal(err)
   }
   var token refresh_token
   err = xml.Unmarshal(data, &token)
   if err != nil {
      t.Fatal(err)
   }
   resp, err := token.explore()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("explore.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
