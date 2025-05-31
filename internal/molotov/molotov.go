package main

import (
   "41.neocities.org/media/molotov"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.password != "":
      err = set.authenticate()
   case set.address.String() != "":
      err = set.download()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) authenticate() error {
   var login molotov.Login
   err := login.New(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := login.Auth.Refresh()
   if err != nil {
      return err
   }
   return write_file(f.media+"/molotov/Refresh", data)
}

func (f *flag_set) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/molotov/Asset")
      if err != nil {
         return err
      }
      var asset molotov.Asset
      err = asset.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return asset.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/molotov/Refresh")
   if err != nil {
      return err
   }
   var refresh molotov.Refresh
   err = refresh.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = refresh.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.media+"/molotov/Refresh", data)
   if err != nil {
      return err
   }
   view, err := refresh.View(&f.address)
   if err != nil {
      return err
   }
   data, err = refresh.Asset(view)
   if err != nil {
      return err
   }
   var asset molotov.Asset
   err = asset.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/molotov/Asset", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(asset.FhdReady())
   if err != nil {
      return err
   }
   return net.Mpd(f.media+"/Mpd", resp)
}

type flag_set struct {
   address  molotov.Address
   dash     string
   e        net.License
   email    string
   media    string
   password string
}
