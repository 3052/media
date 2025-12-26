package disney

import (
   "encoding/xml"
   "log"
   "os"
   "os/exec"
   "testing"
)

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestLogin(t *testing.T) {
   log.SetFlags(log.Ltime)
   email, err := output("credential", "-h=disneyplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=disneyplus.com")
   if err != nil {
      t.Fatal(err)
   }
   var device_item device
   err = device_item.register()
   if err != nil {
      t.Fatal(err)
   }
   account_without, err := device_item.login(email, password)
   if err != nil {
      t.Fatal(err)
   }
   account_with, err := account_without.switch_profile()
   if err != nil {
      t.Fatal(err)
   }
   data, err := xml.Marshal(account_with)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = write_file(cache+"/disney/account.xml", data)
   if err != nil {
      t.Fatal(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
