package main

import (
   "41.neocities.org/media/criterion"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
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
      err = set.do_token()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/criterion/Token")
   if err != nil {
      return err
   }
   var token criterion.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = token.Refresh()
   if err != nil {
      return err
   }
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/criterion/Token", data)
   if err != nil {
      return err
   }
   video, err := token.Video(path.Base(f.address))
   if err != nil {
      return err
   }
   data, err = token.Files(video)
   if err != nil {
      return err
   }
   var files criterion.Files
   err = files.Unmarshal(data)
   if err != nil {
      return err
   }
   file, ok := files.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(file.Links.Source.Href)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return file.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_token() error {
   data, err := criterion.NewToken(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/criterion/Token", data)
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

type flag_set struct {
   address  string
   config   net.Config
   email    string
   filters  net.Filters
   cache    string
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
