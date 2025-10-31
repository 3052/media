package main

import (
   "41.neocities.org/media/hulu"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "net/url"
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
   
   //f.config.ClientId = f.cache + "/L3/client_id.bin"
   //f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   //flag.StringVar(&f.config.ClientId, "C", f.config.ClientId, "client ID")
   //flag.StringVar(&f.config.PrivateKey, "P", f.config.PrivateKey, "private key")
   
   f.config.CertificateChain = f.cache + "/SL2000/CertificateChain"
   f.config.EncryptSignKey = f.cache + "/SL2000/EncryptSignKey"
   flag.StringVar(&f.config.CertificateChain, "C", f.config.CertificateChain, "certificate chain")
   flag.StringVar(&f.config.EncryptSignKey, "E", f.config.EncryptSignKey, "encrypt sign key")
   
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/hulu/Authenticate")
   if err != nil {
      return err
   }
   var auth hulu.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   id, err := hulu.Id(f.address)
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(id)
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   var play hulu.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      //return play.Widevine(data)
      return play.PlayReady(data)
   }
   return f.filters.Filter(resp, &f.config)
}

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         if filepath.Ext(req.URL.Path) != ".mp4" {
            log.Println(req.Method, req.URL)
         }
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
      err = set.do_authenticate()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   email    string
   filters  net.Filters
   password string
}

func (f *flag_set) do_authenticate() error {
   data, err := hulu.NewAuthenticate(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/hulu/Authenticate", data)
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
