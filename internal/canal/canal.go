package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return cache.Player.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

func (f *flag_set) do_episode_movie() error {
   data, err := os.ReadFile(f.cache + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Player, err = cache.Session.Player(f.tracking_id)
   if err != nil {
      return err
   }
   err = cache.Player.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/canal/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
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
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   flag.Int64Var(&f.season, "s", 0, "season")
   flag.StringVar(&f.tracking_id, "t", "", "tracking ID")
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
      if path.Ext(req.URL.Path) == ".dash" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email_password() {
      err = set.do_session()
   } else if set.refresh {
      err = set.do_refresh()
   } else if set.address != "" {
      err = set.do_address()
   } else if set.tracking_id != "" {
      if set.season >= 1 {
         err = set.do_season()
      } else {
         err = set.do_episode_movie()
      }
   } else if set.dash != "" {
      err = set.do_dash()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_session() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   token, err := ticket.Token(f.email, f.password)
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Fetch(token.SsoToken)
   if err != nil {
      return err
   }
   data, err := json.Marshal(canal.Cache{Session: &session})
   if err != nil {
      return err
   }
   return write_file(f.cache+"/canal/Cache", data)
}

func (f *flag_set) do_refresh() error {
   data, err := os.ReadFile(f.cache + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   err = cache.Session.Fetch(cache.Session.SsoToken)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/canal/Cache", data)
}

func (f *flag_set) do_address() error {
   tracking_id, err := canal.TrackingId(f.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking id =", tracking_id)
   return nil
}

func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   episodes, err := cache.Session.Episodes(f.tracking_id, f.season)
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

type flag_set struct {
   config      net.Config
   cache       string
   email       string
   password    string
   refresh     bool
   address     string
   tracking_id string
   season      int64
   dash        string
}
