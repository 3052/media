package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/paramount"
   "41.neocities.org/platform/proxy"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.paramount, "b", "", "paramount ID")
   flag.StringVar(&f.e.ClientId, "client", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.BoolVar(&f.intl, "intl", false, "P+ instance: INTL")
   flag.StringVar(&f.e.PrivateKey, "key", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.proxy, "p", false, "proxy")
   flag.Parse()
   if f.proxy {
      http.DefaultClient.Transport = &proxy.Transport{
         Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
         Proxy:     http.ProxyFromEnvironment,
      }
   }
   if f.paramount != "" {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
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
      session, err := at.Session(f.paramount)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return session.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
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

type flags struct {
   dash      string
   e         internal.License
   intl      bool
   media     string
   paramount string
   proxy     bool
}
