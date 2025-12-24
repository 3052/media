package disney

import (
   "io"
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := playback()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("disney.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
