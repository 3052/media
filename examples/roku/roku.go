package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/roku"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (r *runner) do_dash() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.Widevine(data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      return "L"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) do_user() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   cache.User, err = cache.Connection.GetUser(cache.LinkCode)
   if err != nil {
      return err
   }
   return r.write(cache)
}

type runner struct {
   cache  string
   config maya.Config
   // 1
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   // 4
   dash string
}

func (r *runner) do_connection() error {
   var (
      cache roku.Cache
      err   error
   )
   cache.Connection, err = roku.NewConnection(nil)
   if err != nil {
      return err
   }
   cache.LinkCode, err = cache.Connection.RequestLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(cache.LinkCode)
   return r.write(&cache)
}

func (r *runner) do_roku() error {
   cache := &roku.Cache{}
   if r.get_user {
      var err error
      cache, err = r.read()
      if err != nil {
         return err
      }
   }
   connection, err := roku.NewConnection(cache.User)
   if err != nil {
      return err
   }
   cache.Playback, err = connection.Playback(r.roku)
   if err != nil {
      return err
   }
   err = cache.Playback.Mpd(cache)
   if err != nil {
      return err
   }
   err = r.write(cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.MpdBody, cache.Mpd)
}

func (r *runner) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   r.cache = cache + "/roku/Cache"
   r.config.ClientId = cache + "/L3/client_id.bin"
   r.config.PrivateKey = cache + "/L3/private_key.pem"

   flag.StringVar(&r.config.ClientId, "C", r.config.ClientId, "client ID")
   flag.StringVar(&r.config.PrivateKey, "P", r.config.PrivateKey, "private key")
   flag.BoolVar(&r.connection, "c", false, "connection")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.BoolVar(&r.get_user, "g", false, "get user")
   flag.StringVar(&r.roku, "r", "", "Roku ID")
   flag.BoolVar(&r.set_user, "s", false, "set user")
   flag.Parse()
   if r.connection {
      return r.do_connection()
   }
   if r.set_user {
      return r.do_user()
   }
   if r.roku != "" {
      return r.do_roku()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

func (r *runner) write(cache *roku.Cache) error {
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", r.cache)
   return os.WriteFile(r.cache, data, os.ModePerm)
}

func (r *runner) read() (*roku.Cache, error) {
   data, err := os.ReadFile(r.cache)
   if err != nil {
      return nil, err
   }
   var cache roku.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return nil, err
   }
   return &cache, nil
}
