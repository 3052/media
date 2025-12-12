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

type user_cache struct {
   Mpd     *url.URL
   MpdBody []byte
   Player  *Player
   Session Session
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/canal/user_cache.xml"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.BoolVar(&c.subtitles, "S", false, "subtitles")
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   flag.Int64Var(&c.season, "s", 0, "season")
   flag.StringVar(&c.tracking, "t", "", "tracking")
   flag.Parse()

   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.tracking != "" {
      if c.season >= 1 {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   if c.subtitles {
      return c.do_subtitles()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *command) do_email_password() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   token, err := ticket.Token(c.email, c.password)
   if err != nil {
      return err
   }
   var cache user_cache
   err = cache.Session.Fetch(token.SsoToken)
   if err != nil {
      return err
   }
   return write(c.name, &cache)
}

type command struct {
   config   net.Config
   name    string
   // 1
   email    string
   password string

   // 2
   refresh  bool
   // 3
   address  string
   // 4
   tracking string
   season   int64
   // 5
   subtitles bool
   // 6
   dash     string
}

///

func (c *command) do_refresh() error {
   data, err := os.ReadFile(c.name + "/canal/user_cache")
   if err != nil {
      return err
   }
   var cache canal.user_cache
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
   return write_file(c.name+"/canal/user_cache", data)
}

func (c *command) do_address() error {
   tracking, err := canal.Tracking(c.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking =", tracking)
   return nil
}

func (c *command) do_tracking_season() error {
   data, err := os.ReadFile(c.name + "/canal/user_cache")
   if err != nil {
      return err
   }
   var cache canal.user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   episodes, err := cache.Session.Episodes(c.tracking, c.season)
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

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name + "/canal/user_cache")
   if err != nil {
      return err
   }
   var cache canal.user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Player.Widevine(data)
   }
   return c.config.Download(cache.MpdBody, cache.Mpd, c.dash)
}

func (c *command) do_tracking() error {
   data, err := os.ReadFile(c.name + "/canal/user_cache")
   if err != nil {
      return err
   }
   var cache canal.user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Player, err = cache.Session.Player(c.tracking)
   if err != nil {
      return err
   }
   if c.subtitles {
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
   err = write_file(c.name+"/canal/user_cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}
