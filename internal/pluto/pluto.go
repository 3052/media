package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/pluto"
   "errors"
   "flag"
   "os"
   "path/filepath"
)

type flags struct {
   address pluto.Address
   dash    string
   e       internal.License
   media   string
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
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&pluto.ForwardedFor, "s", "", "set forward")
   flag.Parse()
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
      f.e.Widevine = pluto.Widevine
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   video, err := f.address.Vod()
   if err != nil {
      return err
   }
   clips, err := video.Clips()
   if err != nil {
      return err
   }
   file, ok := clips.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := file.Mpd()
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}
