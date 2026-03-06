package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/canal"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
)

func (c *client) do_email_password() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.email, c.password)
   if err != nil {
      return err
   }
   c.Session = &canal.Session{}
   err = c.Session.Fetch(login.SsoToken)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_refresh() error {
   return cache.Update(c, func() error {
      return c.Session.Fetch(c.Session.SsoToken)
   })
}

func (c *client) do_address() error {
   tracking, err := canal.FetchTracking(c.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking =", tracking)
   return nil
}

func (c *client) do_tracking() error {
   err := cache.Update(c, func() error {
      var err error
      c.Player, err = c.Session.Player(c.tracking)
      if err != nil {
         return err
      }
      c.Dash, err = c.Player.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_tracking_season() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   episodes, err := c.Session.Episodes(c.tracking, c.season)
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

type client struct {
   Dash    *canal.Dash
   Player  *canal.Player
   Session *canal.Session
   // 1
   email    string
   password string
   // 2
   refresh bool
   // 3
   address string
   // 4
   tracking string
   season   int
   // 5
   subtitles bool
   // 6
   dash_id string
}

func (c *client) do_subtitles() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   for _, subtitles := range c.Player.Subtitles {
      err = get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}
func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Player.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
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
   maya.SetProxy("", "*.dash")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/canal.xml")
   if err != nil {
      return err
   }
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
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
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
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"r"},
      {"a"},
      {"t", "s"},
      {"S"},
      {"d", "C", "P"},
   })
}
