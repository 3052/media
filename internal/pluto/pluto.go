package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "errors"
   "flag"
   "os"
   "path/filepath"
)

///

type flags struct {
   e     net.License
   media string
   
   address string
   dash    string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&pluto.ForwardedFor, "s", "", "set forward")
   flag.IntVar(&net.ThreadCount, "t", 1, "thread count")
   flag.Parse()
   switch {
   case f.address != "":
      err = f.do_address()
   case f.dash != "":
      err = f.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flags) do_address() error {
   var address pluto.Address
   err := address.Set(f.address)
   if err != nil {
      return err
   }
   video, err := address.Vod()
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
   return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   f.e.Widevine = pluto.Widevine
   return f.e.Download(f.media+"/Mpd", f.dash)
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
