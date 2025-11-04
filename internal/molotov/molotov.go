package main

import (
   "41.neocities.org/media/molotov"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/molotov/Refresh")
   if err != nil {
      return err
   }
   var refresh molotov.Refresh
   err = refresh.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = refresh.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/molotov/Refresh", data)
   if err != nil {
      return err
   }
   var address molotov.Address
   err = address.Parse(f.address)
   if err != nil {
      return err
   }
   view, err := refresh.View(&address)
   if err != nil {
      return err
   }
   asset, err := refresh.Asset(view)
   if err != nil {
      return err
   }
   resp, err := http.Get(asset.FhdReady())
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return asset.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_refresh() error {
   var login molotov.Login
   err := login.New(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := login.Auth.Refresh()
   if err != nil {
      return err
   }
   return write_file(f.cache+"/molotov/Refresh", data)
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   email    string
   filters  net.Filters
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
func main() {
   http.DefaultTransport = &molotov.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.address != "":
      err = set.do_address()
   case set.email_password():
      err = set.do_refresh()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}
