package main

import (
   "41.neocities.org/media/ctv"
   "41.neocities.org/net"
   "flag"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   media, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   media = filepath.ToSlash(media) + "/media"
   f.cdm.ClientId = media + "/client_id.bin"
   f.cdm.PrivateKey = media + "/private_key.pem"
   f.filters = net.Filters{
      {BitrateStart: 100_000, BitrateEnd: 200_000},
      {BitrateStart: 2_000_000, BitrateEnd: 7_000_000},
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.IntVar(&net.Threads, "t", 5, "threads")
   flag.Func("x", "proxy", func(data string) error {
      var err error
      f.proxy, err = url.Parse(data)
      return err
   })
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   http.DefaultTransport = net.Transport(f.proxy)
   if set.address != "" {
      err = set.do_address()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flag_set) do_address() error {
   resolve, err := ctv.Resolve(ctv.Path(f.address))
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
   f.cdm.License = ctv.License
   return f.filters.Filter(resp, &f.cdm)
}

type flag_set struct {
   address string
   cdm     net.Cdm
   filters net.Filters
   proxy   *url.URL
}

