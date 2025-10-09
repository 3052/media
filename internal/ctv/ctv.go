package main

import (
   "41.neocities.org/media/ctv"
   "41.neocities.org/net"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

type flag_set struct {
   address string
   cdm     net.Cdm
   filters net.Filters
}

func (f *flag_set) New() error {
   media, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   media = filepath.ToSlash(media) + "/media"
   f.cdm.ClientId = media + "/client_id.bin"
   f.cdm.PrivateKey = media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
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
