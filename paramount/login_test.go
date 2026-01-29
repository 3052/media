package paramount

import (
   "log"
   "os"
   "os/exec"
   "testing"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func TestLogin(t *testing.T) {
   user, err := output("credential", "-h=paramountplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=paramountplus.com")
   if err != nil {
      t.Fatal(err)
   }
   at, err := GetAt(ComCbsApp.AppSecret)
   if err != nil {
      t.Fatal(err)
   }
   cookie, err := Login(at, user, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = write_file(
      cache + "/paramount/login.txt", []byte(cookie.String()),
   )
   if err != nil {
      t.Fatal(err)
   }
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}
