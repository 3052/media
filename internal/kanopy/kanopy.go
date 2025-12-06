package main

import (
   "41.neocities.org/media/kanopy"
   "41.neocities.org/net"
   "encoding/json"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (r *runner) run() error {
   var err error
   r.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   r.cache = filepath.ToSlash(r.cache)
   r.config.ClientId = r.cache + "/L3/client_id.bin"
   r.config.PrivateKey = r.cache + "/L3/private_key.pem"
   flag.StringVar(&r.config.ClientId, "C", r.config.ClientId, "client ID")
   flag.StringVar(&r.config.PrivateKey, "P", r.config.PrivateKey, "private key")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.StringVar(&r.email, "e", "", "email")
   flag.IntVar(&r.kanopy, "k", 0, "Kanopy ID")
   flag.StringVar(&r.password, "p", "", "password")
   flag.IntVar(&r.config.Threads, "t", 2, "threads")
   flag.Parse()
   if r.email != "" {
      if r.password != "" {
         return r.do_login()
      }
   }
   if r.kanopy >= 1 {
      return r.do_kanopy()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

func (r *runner) do_kanopy() error {
   data, err := os.ReadFile(r.cache + "/kanopy/Cache")
   if err != nil {
      return err
   }
   var cache kanopy.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   member, err := cache.Login.Membership()
   if err != nil {
      return err
   }
   plays, err := cache.Login.Plays(member, r.kanopy)
   if err != nil {
      return err
   }
   var ok bool
   cache.Manifest, ok = plays.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   err = cache.Manifest.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(r.cache + "/kanopy/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (r *runner) do_dash() error {
   data, err := os.ReadFile(r.cache + "/kanopy/Cache")
   if err != nil {
      return err
   }
   var cache kanopy.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return cache.Login.Widevine(cache.Manifest, data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   var tool runner
   err := tool.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) do_login() error {
   var login kanopy.Login
   err := login.Fetch(r.email, r.password)
   if err != nil {
      return err
   }
   data, err := json.Marshal(kanopy.Cache{Login: &login})
   if err != nil {
      return err
   }
   return write_file(r.cache+"/kanopy/Cache", data)
}

type runner struct {
   cache    string
   config   net.Config
   // 1
   email    string
   password string
   // 2
   kanopy   int
   // 3
   dash string
}
