package main

import (
   "41.neocities.org/media/molotov"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   return nil
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media+name)
   return os.WriteFile(f.media+name, data, os.ModePerm)
}

type flags struct {
   e        internal.License
   email    string
   password string
   media    string
   dash     string
   address molotov.Address
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Var(&f.address, "a", "address")
   flag.Parse()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.address.String() != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) authenticate() error {
   var login molotov.Login
   err := login.New(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := login.Auth.Refresh()
   if err != nil {
      return err
   }
   return f.write_file("/molotov/Refresh", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      
      data, err := os.ReadFile(f.media + "/molotov/Playlist")
      if err != nil {
         return err
      }
      var play molotov.Playlist
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
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
   err = f.write_file("/molotov/Refresh", data)
   if err != nil {
      return err
   }
   view, err := refresh.View(&f.address)
   if err != nil {
      return err
   }
   asset, err := refresh.Asset(view)
   if err != nil {
      return err
   }
   resp, err := http.Get(asset.FhdReady())
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}
