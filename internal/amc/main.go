package main

import (
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
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4f" {
         return ""
      }
      return "LP"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) run() error {
   var err error
   r.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   r.cache = filepath.ToSlash(r.cache)
   r.config.ClientId = r.cache + "/L3/client_id.bin"
   r.config.PrivateKey = r.cache + "/L3/private_key.pem"

   flag.StringVar(&r.email, "E", "", "email")
   flag.StringVar(&r.password, "P", "", "password")
   flag.Int64Var(&r.series, "S", 0, "series ID")
   flag.StringVar(&r.config.ClientId, "c", r.config.ClientId, "client ID")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.Int64Var(&r.episode, "e", 0, "episode or movie ID")
   flag.StringVar(&r.config.PrivateKey, "p", r.config.PrivateKey, "private key")
   flag.BoolVar(&r.refresh, "r", false, "refresh")
   flag.Int64Var(&r.season, "s", 0, "season ID")
   flag.Parse()

   if r.email != "" {
      if r.password != "" {
         return r.do_auth()
      }
   }
   if r.refresh {
      return r.do_refresh()
   }
   if r.series >= 1 {
      return r.do_series()
   }
   if r.season >= 1 {
      return r.do_season()
   }
   if r.episode >= 1 {
      return r.do_episode()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

type runner struct {
   config   net.Config
   dash     string
   email    string
   episode  int64
   password string
   refresh  bool
   season   int64
   series   int64
   cache string
}
