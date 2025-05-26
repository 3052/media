package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "os"
   "path/filepath"
)

func (f *flags) do_paramount() error {
   paramount.Client.Transport = net.Proxy(false)
   secret := paramount.ComCbsApp
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   at, err := secret.At()
   if err != nil {
      return err
   }
   session, err := at.Session(f.paramount)
   if err != nil {
      return err
   }
   f.license.Widevine = func(data []byte) ([]byte, error) {
      return session.Widevine(data)
   }
   if f.intl {
      secret = paramount.ComCbsCa
   }
   at, err = secret.At()
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
   return f.license.Bitrate(resp, &f.bitrate)
}

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.license.ClientId = f.media + "/client_id.bin"
   f.license.PrivateKey = f.media + "/private_key.pem"
   f.bitrate.Value = [][2]int{
      {100_000, 150_000}, {3_900_000, 5_900_000},
   }
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   flag.BoolVar(&f.intl, "i", false, "P+ instance: INTL")
   flag.StringVar(
      &f.license.PrivateKey, "k", f.license.PrivateKey, "private key",
   )
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Parse()
   if f.paramount != "" {
      err = f.do_paramount()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

type flags struct {
   media string
   license     net.License
   ///////////////////////
   paramount string
   intl      bool
   bitrate net.Bitrate
}
