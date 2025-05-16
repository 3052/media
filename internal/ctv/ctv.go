package main

import (
   "41.neocities.org/media/ctv"
   "41.neocities.org/stream"
   "flag"
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

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Func("a", "address", func(data string) error {
      return f.address.Set(data)
   })
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
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

type flags struct {
   e     stream.License
   media string

   address ctv.Address
   dash    string
}

///

func (f *flags) do_address() error {
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
   return stream.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
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
   return stream.Mpd(f.media+"/Mpd", resp)
}
