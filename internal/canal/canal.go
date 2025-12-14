package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/canal"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_email_password() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.email, c.password)
   if err != nil {
      return err
   }
   var cache user_cache
   err = cache.Session.Fetch(login.SsoToken)
   if err != nil {
      return err
   }
   return write(c.name, &cache)
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func (c *command) do_refresh() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   err = cache.Session.Fetch(cache.Session.SsoToken)
   if err != nil {
      return err
   }
   return write(c.name, cache)
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
   cache, err := read(c.name)
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

func (c *command) do_tracking() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.Player, err = cache.Session.Player(c.tracking)
   if err != nil {
      return err
   }
   cache.Mpd, cache.MpdBody, err = cache.Player.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd, cache.MpdBody)
}

func (c *command) do_subtitles() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   for _, subtitles := range cache.Player.Subtitles {
      err = get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

type command struct {
   address   string
   config    maya.Config
   dash      string
   email     string
   name      string
   password  string
   refresh   bool
   season    int64
   subtitles bool
   tracking  string
}
type user_cache struct {
   Mpd     *url.URL
   MpdBody []byte
   Player  *canal.Player
   Session canal.Session
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Player.Widevine(data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
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

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
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
