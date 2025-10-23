package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
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

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
      Proxy: func(req *http.Request) (*url.URL, error) {
         switch {
         case filepath.Base(req.URL.Path) == "anonymous-session-token.json":
            log.Println(req.Method, req.URL)
            return nil, nil
         case filepath.Base(req.URL.Path) == "getlicense":
            log.Println(req.Method, req.URL)
            return nil, nil
         case filepath.Ext(req.URL.Path) == ".m4s":
            return nil, nil
         }
         log.Println(req.Method, req.URL)
         return http.ProxyFromEnvironment(req)
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.paramount != "" {
      err := set.do_paramount()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flag_set) secret() paramount.AppSecret {
   if f.intl {
      return paramount.ComCbsCa
   }
   return paramount.ComCbsApp
}

type flag_set struct {
   cache     string
   config    net.Config
   filters   net.Filters
   intl      bool
   paramount string
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
