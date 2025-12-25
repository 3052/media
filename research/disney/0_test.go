package disney

import (
   "io"
   "os"
   "testing"
)

func TestRefreshToken(t *testing.T) {
   resp, err := refresh_token()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("refresh_token.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
