package main

import (
   "41.neocities.org/media/molotov"
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
      Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
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
   func() {
      if set.email != "" {
         if set.password != "" {
            err = set.authenticate()
            return
         }
      }
      if set.address.String() != "" {
         err = set.download()
      } else {
         flag.Usage()
      }
   }()
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   address  molotov.Address
   dash     string
   cdm      net.Cdm
   filters  net.Filters
   email    string
   cache    string
   password string
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) authenticate() error {
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

func (f *flag_set) download() error {
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
   view, err := refresh.View(&f.address)
   if err != nil {
      return err
   }
   data, err = refresh.Asset(view)
   if err != nil {
      return err
   }
   var asset molotov.Asset
   err = asset.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(asset.FhdReady())
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return asset.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.cdm.ClientId = f.cache + "/L3/client_id.bin"
   f.cdm.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.cdm.PrivateKey, "k", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}
