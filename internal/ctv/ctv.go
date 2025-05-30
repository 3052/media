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

func (f *flag_set) do_address() error {
   http.DefaultTransport = ctv.Transport(f.proxy)
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
   f.cdm.License = ctv.License
   return f.filters.Filter(resp, &f.cdm)
}

func main() {
   var f flag_set
   err := f.New()
   if err != nil {
      panic(err)
   }
   if f.address != "" {
      err = f.do_address()
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
   ///////////////////////////////////////////////////
   flag.Func("a", "address", func(data string) error {
      return f.address.Set(data)
   })
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.Func("x", "proxy", func(data string) error {
      var err error
      f.proxy, err = url.Parse(data)
      return err
   })
   ///////////////////////////////////////////////////////////////////////
   flag.Parse()
   return nil
}

type flag_set struct {
   address ctv.Address
   cdm     net.Cdm
   filters net.Filters
   proxy *url.URL
}
