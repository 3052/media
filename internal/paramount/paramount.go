package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
   "strings"
)

func (f *flag_set) do_paramount() error {
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   at, err := paramount.ComCbsApp.At()
   if err != nil {
      return err
   }
   data, err := at.Session(f.paramount)
   if err != nil {
      return err
   }
   var session paramount.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   at, err = f.secret().At()
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
   f.config.Send = func(data []byte) ([]byte, error) {
      return session.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

func (f *flag_set) secret() paramount.AppSecret {
   if f.intl {
      return paramount.ComCbsCa
   }
   return paramount.ComCbsApp
}

type flag_set struct {
   cache      string
   config     net.Config
   filters    net.Filters
   intl       bool
   paramount  string
   bypass string // Added field
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.config.ClientId, "C", f.config.ClientId, "client ID")
   flag.StringVar(&f.config.PrivateKey, "P", f.config.PrivateKey, "private key")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.BoolVar(&f.intl, "i", false, "intl")
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.StringVar(&f.bypass, "b", ".m4s,.mp4", "proxy bypass")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   log.SetFlags(log.Ltime)
   http.DefaultTransport = net.Proxy(func(req *http.Request) bool {
      for _, ext := range strings.Split(set.bypass, ",") {
         if filepath.Ext(req.URL.Path) == ext {
            return true
         }
      }
      return false
   })
   if set.paramount != "" {
      err := set.do_paramount()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
