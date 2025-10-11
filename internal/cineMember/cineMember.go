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
   email    string
   filters  net.Filters
   password string
}

func (f *flag_set) do_user() error {
   data, err := cineMember.NewUser(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/cineMember/User", data)
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
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/cineMember/User")
   if err != nil {
      return err
   }
   var user cineMember.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   var address cineMember.Address
   err = address.Parse(f.address)
   if err != nil {
      return err
   }
   article, err := address.Article()
   if err != nil {
      return err
   }
   asset, ok := article.Film()
   if !ok {
      return errors.New(".Film()")
   }
   data, err = user.Play(article, asset)
   if err != nil {
      return err
   }
   var play cineMember.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   title, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(title.Manifest)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return title.Send(data)
   }
   return f.filters.Filter(resp, &f.config)
}
