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

type flags struct {
   content_id string
   dash       string
   e          internal.License
   media      string
   mullvad    bool
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.content_id, "b", "", "content ID")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
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

func (f *flags) download() error {
   if f.dash != "" {
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
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return session.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   var secret paramount.AppSecret
   if f.mullvad {
      secret = paramount.ComCbsCa
      http.DefaultClient.Transport = &mullvad.Transport{
         Proxy: http.ProxyFromEnvironment,
      }
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
   return internal.Mpd(f.media+"/Mpd", resp)
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
