package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/canal"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Player.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash     *canal.Dash
   Player  *canal.Player
   Session *canal.Session
}

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
   var session canal.Session
   err = session.Fetch(login.SsoToken)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Session: &session})
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

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/canal/userCache.xml"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.tracking, "t", "", "tracking")
   flag.IntVar(&c.season, "s", 0, "season")
   // 5
   flag.BoolVar(&c.subtitles, "S", false, "subtitles")
   // 6
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.Parse()
   // 1
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   // 2
   if c.refresh {
      return c.do_refresh()
   }
   // 3
   if c.address != "" {
      return c.do_address()
   }
   // 4
   if c.tracking != "" {
      if c.season >= 1 {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   // 5
   if c.subtitles {
      return c.do_subtitles()
   }
   // 6
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
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
   tracking, err := canal.FetchTracking(c.address)
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
   job    maya.WidevineJob
   name      string
   // 1
   email     string
   password  string
   // 2
   refresh   bool
   // 3
   address   string
   // 4
   tracking  string
   season    int
   // 5
   subtitles bool
   // 6
   dash      string
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
   cache.Dash, err = cache.Player.Dash()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}
