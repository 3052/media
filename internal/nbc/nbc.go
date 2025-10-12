package main

import (
   "41.neocities.org/media/nbc"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return http.ProxyFromEnvironment(req)
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.nbc >= 1 {
      err = set.do_nbc()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

type flag_set struct {
   cache   string
   config  net.Config
   filters net.Filters
   nbc     int
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
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.IntVar(&f.nbc, "n", 0, "NBC ID")
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.Parse()
   return nil
}

func (f *flag_set) do_nbc() error {
   var metadata nbc.Metadata
   err := metadata.New(f.nbc)
   if err != nil {
      return err
   }
   vod, err := metadata.Vod()
   if err != nil {
      return err
   }
   resp, err := http.Get(vod.PlaybackUrl)
   if err != nil {
      return err
   }
   f.config.Send = nbc.Widevine
   return f.filters.Filter(resp, &f.config)
}
