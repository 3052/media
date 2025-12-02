package main

import (
   "41.neocities.org/media/hulu"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (f *flag_set) do_session() error {
   var session hulu.Session
   err := session.Fetch(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := json.Marshal(hulu.Cache{Session: &session})
   if err != nil {
      return err
   }
   return write_file(f.cache+"/hulu/Cache", data)
}

func main() {
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
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
      err = set.do_session()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   dash string
   email    string
   password string
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
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/hulu/Cache")
   if err != nil {
      return err
   }
   var session hulu.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   err = session.Refresh()
   if err != nil {
      return err
   }
   id, err := hulu.Id(f.address)
   if err != nil {
      return err
   }
   deep, err := session.DeepLink(id)
   if err != nil {
      return err
   }
   playlist, err := session.Playlist(deep)
   if err != nil {
      return err
   }
   resp, err := http.Get(playlist.StreamUrl)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return playlist.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
