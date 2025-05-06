package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/paramount"
   "flag"
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
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.BoolVar(&f.intl, "i", false, "P+ instance: INTL")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.Parse()
   switch {
   case f.paramount != "":
      err = f.do_paramount()
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
   media     string
   e         internal.License
   
   paramount string
   intl      bool
   
   dash      string
}

func (f *flags) do_paramount() error {
   var secret paramount.AppSecret
   if f.intl {
      secret = paramount.ComCbsCa
   } else {
      secret = paramount.ComCbsApp
   }
   at, err := secret.At()
   if err != nil {
      return err
   }
   item, err := at.Item(f.paramount)
   if err != nil {
      return err
   }
   resp, err := item.Mpd()
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   at, err := paramount.ComCbsApp.At()
   if err != nil {
      return err
   }
   session, err := at.Session(f.paramount)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return session.Widevine(data)
   }
   return f.e.Download(f.media+"/Mpd", f.dash)
}
