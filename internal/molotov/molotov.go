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
   media    string
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
   return write_file(f.media+"/molotov/Refresh", data)
}

func (f *flag_set) download() error {
   data, err := os.ReadFile(f.media + "/molotov/Refresh")
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
   err = write_file(f.media+"/molotov/Refresh", data)
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
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.cdm.PrivateKey, "k", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}
