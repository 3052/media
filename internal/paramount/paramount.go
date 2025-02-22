package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/paramount"
   "41.neocities.org/platform/mullvad"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) New() error {
   var err error
   f.home, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.home = filepath.ToSlash(f.home) + "/media"
   f.e.ClientId = f.home + "/client_id.bin"
   f.e.PrivateKey = f.home + "/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.content_id, "b", "", "content ID")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.BoolVar(&f.mullvad, "m", false, "Mullvad")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.Parse()
   switch {
   case f.content_id != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

type flags struct {
   content_id     string
   e              internal.License
   home           string
   mullvad        bool
   representation string
}

func (f *flags) download() error {
   if f.representation != "" {
      // INTL does NOT allow anonymous key request, so if you are INTL you
      // will need to use US VPN until someone codes the INTL login
      at, err := paramount.ComCbsApp.At()
      if err != nil {
         return err
      }
      session, err := at.Session(f.content_id)
      if err != nil {
         return err
      }
      f.e.Widevine = session.Widevine()
      return f.e.Download(f.home, f.representation)
   }
   var secret paramount.AppSecret
   if f.mullvad {
      secret = paramount.ComCbsCa
      http.DefaultClient.Transport = new(mullvad.Vpn)
      defer mullvad.Disconnect()
   } else {
      secret = paramount.ComCbsApp
   }
   at, err := secret.At()
   if err != nil {
      return err
   }
   item, err := at.Item(f.content_id)
   if err != nil {
      return err
   }
   resp, err := item.Mpd()
   if err != nil {
      return err
   }
   return internal.Mpd(resp, f.home)
}
