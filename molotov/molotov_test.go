package molotov

import (
   "41.neocities.org/widevine"
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestWidevine(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/molotov/refresh")
   if err != nil {
      t.Fatal(err)
   }
   var token refresh
   err = token.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   view1, err := token.view()
   if err != nil {
      t.Fatal(err)
   }
   assets1, err := token.assets(view1)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", assets1.fhd_ready())
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   pssh.KeyIds = [][]byte{
      []byte("\xc3\x1c\xd0+m\x17\x01\xee\xa1\xedp7\xa8~\xd8J"),
   }
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = cdm.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = assets1.widevine(data)
   if err != nil {
      t.Fatal(err)
   }
}

func TestRefresh(t *testing.T) {
   data, err := exec.Command("password", "molotov.tv").Output()
   if err != nil {
      t.Fatal(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   var login1 login
   err = login1.New(email, password)
   if err != nil {
      t.Fatal(err)
   }
   data, err = login1.Auth.refresh()
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/molotov/refresh", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
