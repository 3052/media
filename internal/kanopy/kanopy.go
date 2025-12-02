package main

import (
   "41.neocities.org/media/kanopy"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

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
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.IntVar(&f.kanopy, "k", 0, "Kanopy ID")
   flag.StringVar(&f.password, "p", "", "password")
   flag.IntVar(&f.config.Threads, "t", 2, "threads")
   flag.Parse()
   return nil
}

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = net.Proxy(func(req *http.Request) bool {
      return filepath.Ext(req.URL.Path) == ".m4s"
   })
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.email_password():
      err = set.do_login()
   case set.kanopy >= 1:
      err = set.do_kanopy()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func (f *flag_set) do_login() error {
   data, err := kanopy.FetchLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/kanopy/Login", data)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_kanopy() error {
   data, err := os.ReadFile(f.cache + "/kanopy/Login")
   if err != nil {
      return err
   }
   var login kanopy.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   member, err := login.Membership()
   if err != nil {
      return err
   }
   plays, err := login.Plays(member, f.kanopy)
   if err != nil {
      return err
   }
   manifest, ok := plays.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := manifest.Mpd()
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return login.Widevine(manifest, data)
   }
   return f.filters.Filter(resp, &f.config)
}

type flag_set struct {
   cache    string
   config   net.Config
   email    string
   filters  net.Filters
   kanopy   int
   password string
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}
