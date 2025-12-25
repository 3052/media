package disney

import (
   "encoding/xml"
   "log"
   "os"
   "testing"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func TestRefreshToken(t *testing.T) {
   log.SetFlags(log.Ltime)
   var token refresh_token
   err := token.fetch()
   if err != nil {
      t.Fatal(err)
   }
   data, err := xml.Marshal(token)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = write_file(cache + "/disney/refresh_token.xml", data)
   if err != nil {
      t.Fatal(err)
   }
}
