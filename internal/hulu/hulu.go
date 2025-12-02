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
   flag.IntVar(&f.config.Threads, "t", 2, "threads")
   flag.Parse()
   return nil
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
   case set.email_password():
      err = set.do_session()
   case set.address != "":
      err = set.do_address()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

type flag_set struct {
   cache    string
   config   net.Config
   // 1
   email    string
   password string
   // 2
   address  string
   // 3
   dash     string
}

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

func (f *flag_set) do_address() error {
   id, err := hulu.Id(f.address)
   if err != nil {
      return err
   }
   data, err := os.ReadFile(f.cache + "/hulu/Cache")
   if err != nil {
      return err
   }
   var cache hulu.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   err = cache.Session.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := cache.Session.DeepLink(id)
   if err != nil {
      return err
   }
   cache.Playlist, err = cache.Session.Playlist(deep_link)
   if err != nil {
      return err
   }
   err = cache.Playlist.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(f.cache + "/hulu/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/hulu/Cache")
   if err != nil {
      return err
   }
   var cache hulu.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playlist.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}
