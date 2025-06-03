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
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   http.DefaultTransport = net.Transport(set.proxy)
   if set.paramount != "" {
      err = set.do_paramount()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flag_set) New() error {
   media, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   media = filepath.ToSlash(media) + "/media"
   f.cdm.ClientId = media + "/client_id.bin"
   f.cdm.PrivateKey = media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Func("x", "proxy", func(data string) error {
      var err error
      f.proxy, err = url.Parse(data)
      return err
   })
   f.filters = net.Filters{
      {BitrateStart: 3_000_000, BitrateEnd: 5_000_000},
      {BitrateStart: 100_000, BitrateEnd: 150_000, Role: "main"},
   }
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.Parse()
   return nil
}

type flag_set struct {
   cdm       net.Cdm
   filters   net.Filters
   paramount string
   proxy     *url.URL
}

func (f *flag_set) do_paramount() error {
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
   f.cdm.License = func(data []byte) ([]byte, error) {
      return session.License(data)
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
   return f.filters.Filter(resp, &f.cdm)
}

