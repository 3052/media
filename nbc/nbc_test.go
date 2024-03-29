package nbc

import (
   "154.pages.dev/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "testing"
   "time"
)

func TestVideo(t *testing.T) {
   v, ok := Core().RequestUrl()
   fmt.Println(v, ok)
}

func TestOnDemand(t *testing.T) {
   for _, mpx_guid := range mpx_guids {
      var meta Metadata
      err := meta.New(mpx_guid)
      if err != nil {
         t.Fatal(err)
      }
      video, err := meta.OnDemand()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", video)
      time.Sleep(time.Second)
   }
}

const raw_pssh = "AAAAV3Bzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAADcIARIQBVLkSEJlSk6BsyYAS+R74BoLYnV5ZHJta2V5b3MiEAVS5EhCZUpOgbMmAEvke+AqAkhE"

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   var protect widevine.PSSH
   {
      b, err := base64.StdEncoding.DecodeString(raw_pssh)
      if err != nil {
         t.Fatal(err)
      }
      if err := protect.New(b); err != nil {
         t.Fatal(err)
      }
   }
   module, err := protect.CDM(private_key, client_id)
   if err != nil {
      t.Fatal(err)
   }
   license, err := module.License(Core())
   if err != nil {
      t.Fatal(err)
   }
   key, ok := module.Key(license)
   fmt.Printf("%x %v\n", key, ok)
}
