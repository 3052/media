package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return ""
      }
      switch path.Base(req.URL.Path) {
      case "anonymous-session-token.json", "getlicense":
         return "L"
      }
      return "LP"
   })
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   if set.paramount != "" {
      err := set.do_paramount()
      if err != nil {
         log.Fatal(err)
      }
   } else {
      flag.Usage()
   }
}

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
   flag.Parse()
   return nil
}
