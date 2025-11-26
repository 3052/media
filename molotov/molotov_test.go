package molotov

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/molotov/Login")
   if err != nil {
      t.Fatal(err)
   }
   var login_var Login
   err = login_var.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   var media MediaId
   err = media.Parse("molotov.tv/fr_fr/p/15082-531")
   if err != nil {
      t.Fatal(err)
   }
   _, err = login_var.PlayUrl(&media)
   if err != nil {
      t.Fatal(err)
   }
}
