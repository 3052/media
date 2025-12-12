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

type Cache struct {
   Mpd     *url.URL
   MpdBody []byte
   Player  *Player
   Session *Session
}

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".dash" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   address  string
   name    string
   config   net.Config
   dash     string
   email    string
   password string
   refresh  bool
   season   int64
   subtitles bool
   tracking string
}

///

func (r *command) do_session() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   token, err := ticket.Token(r.email, r.password)
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
   return write_file(r.name+"/canal/Cache", data)
}

func (r *command) do_refresh() error {
   data, err := os.ReadFile(r.name + "/canal/Cache")
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
   return write_file(r.name+"/canal/Cache", data)
}

func (r *command) do_address() error {
   tracking, err := canal.Tracking(r.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking =", tracking)
   return nil
}

func (r *command) do_season() error {
   data, err := os.ReadFile(r.name + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   episodes, err := cache.Session.Episodes(r.tracking, r.season)
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
func (r *command) do_dash() error {
   data, err := os.ReadFile(r.name + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return cache.Player.Widevine(data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}
func (r *command) do_episode_movie() error {
   data, err := os.ReadFile(r.name + "/canal/Cache")
   if err != nil {
      return err
   }
   var cache canal.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Player, err = cache.Session.Player(r.tracking)
   if err != nil {
      return err
   }
   if r.subtitles {
      for _, subtitles := range cache.Player.Subtitles {
         err = get(subtitles.Url)
         if err != nil {
            return err
         }
      }
      return nil
   }
   err = cache.Player.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(r.name+"/canal/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func get(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   file, err := os.Create(path.Base(address))
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   if err != nil {
      return err
   }
   return nil
}

func (r *command) run() error {
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
   flag.BoolVar(&r.subtitles, "S", false, "subtitles")
   flag.StringVar(&r.address, "a", "", "address")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.StringVar(&r.email, "e", "", "email")
   flag.StringVar(&r.password, "p", "", "password")
   flag.BoolVar(&r.refresh, "r", false, "refresh")
   flag.Int64Var(&r.season, "s", 0, "season")
   flag.StringVar(&r.tracking, "t", "", "tracking")
   flag.Parse()
   if r.email != "" {
      if r.password != "" {
         return r.do_session()
      }
   }
   if r.refresh {
      return r.do_refresh()
   }
   if r.address != "" {
      return r.do_address()
   }
   if r.tracking != "" {
      if r.season >= 1 {
         return r.do_season()
      }
      return r.do_episode_movie()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}
