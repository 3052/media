package main

import (
   "41.neocities.org/media/molotov"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
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
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.IntVar(&f.config.Threads, "t", 9, "threads")
   flag.Parse()
   return nil
}

func main() {
   http.DefaultTransport = net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
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
      log.Fatal(err)
   }
}

func (f *flag_set) do_refresh() error {
   login, err := molotov.FetchLogin(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := login.Refresh()
   if err != nil {
      return err
   }
   return write_file(f.cache+"/molotov/Login", data)
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

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/molotov/Login")
   if err != nil {
      return err
   }
   var login molotov.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/molotov/Login", data)
   if err != nil {
      return err
   }
   var media molotov.MediaId
   err = media.Parse(f.address)
   if err != nil {
      return err
   }
   play_url, err := login.PlayUrl(&media)
   if err != nil {
      return err
   }
   playback, err := login.Playback(play_url)
   if err != nil {
      return err
   }
   resp, err := http.Get(playback.FhdReady())
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return playback.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
