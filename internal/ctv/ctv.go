package main

import (
   "41.neocities.org/media/ctv"
   "41.neocities.org/media/internal"
   "41.neocities.org/platform/mullvad"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   address ctv.Address
   dash    string
   e       internal.License
   media   string
   mullvad bool
}

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

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.mullvad, "m", false, "Mullvad")
   flag.Parse()
   if f.mullvad {
      http.DefaultClient.Transport = &mullvad.Transport{}
      defer mullvad.Disconnect()
   }
   switch {
   case f.address[0] != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) download() error {
   if f.dash != "" {
      f.e.Widevine = ctv.Widevine
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   resolve, err := f.address.Resolve()
   if err != nil {
      return err
   }
   axis, err := resolve.Axis()
   if err != nil {
      return err
   }
   content, err := axis.Content()
   if err != nil {
      return err
   }
   address, err := axis.Mpd(content)
   if err != nil {
      return err
   }
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}
