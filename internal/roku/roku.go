package main

import (
   "41.neocities.org/media/roku"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      return "L"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

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
   log.Println("WriteFile", r.cache + "/roku/Cache")
   return os.WriteFile(r.cache + "/roku/Cache", data, os.ModePerm)
}

func (r *runner) read() (*roku.Cache, error) {
   data, err := os.ReadFile(r.cache + "/roku/Cache")
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

func (r *runner) do_connection() error {
   var (
      cache roku.Cache
      user *roku.user
      err error
   )
   cache.Connection, err = user.NewConnection()
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
   cache       string
   config      net.Config
   // 1
   connection  bool
   // 2
   set_user bool
   // 3
   roku        string
   get_user  bool
   // 4
   dash string
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
    connection, err := cache.User.NewConnection()
    if err != nil {
        return err
    }
    playback, err := connection.Playback(r.roku)
    if err != nil {
        return err
    }
    err = playback.Mpd(cache)
    if err != nil {
        return err
    }
    err = r.write(cache)
    if err != nil {
        return err
    }
    return net.Representations(cache.MpdBody, cache.Mpd)
}

func (r *runner) do_dash() error {
   r.config.Send = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return r.filters.Filter(resp, &r.config)
}
