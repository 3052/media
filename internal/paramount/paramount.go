package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   flag.StringVar(
      &f.license.PrivateKey, "k", f.license.PrivateKey, "private key",
   )
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Func("x", "proxy", func(data string) error {
      var err error
      f.proxy, err = url.Parse(data)
      return err
   })
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
   bitrate   net.Bitrate
   license   net.License
   paramount string
   proxy     *url.URL
}

func (f *flags) do_paramount() error {
   http.DefaultTransport = paramount.Transport(f.proxy)
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   secret := paramount.ComCbsApp
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
   if f.proxy != nil {
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
   media, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   media = filepath.ToSlash(media) + "/media"
   f.license.ClientId = media + "/client_id.bin"
   f.license.PrivateKey = media + "/private_key.pem"
   f.bitrate.Value = [][2]int{
      {100_000, 150_000}, {3_900_000, 5_900_000},
   }
   return nil
}
