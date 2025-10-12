package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   filters  net.Filters
   email    string
   password string
}

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
   switch {
   case set.address != "":
      err = set.do_address()
   case set.email_password():
      err = set.do_user()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_user() error {
   var session cineMember.Session
   err := session.New()
   if err != nil {
      return err
   }
   err = session.Login(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(
      f.cache + "/cineMember/Session", []byte(session.String()),
   )
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/cineMember/Session")
   if err != nil {
      return err
   }
   var session cineMember.Session
   err = session.Set(string(data))
   if err != nil {
      return err
   }
   id, err := cineMember.Id(f.address)
   if err != nil {
      return err
   }
   stream, err := session.Stream(id)
   if err != nil {
      return err
   }
   mpd, ok := stream.Mpd()
   if !ok {
      return errors.New(".Mpd()")
   }
   resp, err := http.Get(mpd)
   if err != nil {
      return err
   }
   return f.filters.Filter(resp, &f.config)
}
