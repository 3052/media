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

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.dash")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_tracking() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      var err error
      state.Player, err = state.Session.Player(c.tracking)
      if err != nil {
         return err
      }
      state.Dash, err = state.Player.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_refresh() error {
   var state saved_state
   return c.cache.Update(&state, func() error {
      return state.Session.Fetch(state.Session.SsoToken)
   })
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

func (c *client) do_address() error {
   tracking, err := canal.FetchTracking(c.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking =", tracking)
   return nil
}

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
   var session canal.Session
   err = session.Fetch(login.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Session: &session})
}

type client struct {
   cache maya.Cache
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
   dash string
   job  maya.WidevineJob
}

type saved_state struct {
   Dash    *canal.Dash
   Player  *canal.Player
   Session *canal.Session
}

func (c *client) do_subtitles() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   for _, subtitles := range state.Player.Subtitles {
      err = get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

func (c *client) do_tracking_season() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   episodes, err := state.Session.Episodes(c.tracking, c.season)
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

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/canal.xml")
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
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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
   return maya.Usage([][]string{
      {"e", "p"},
      {"r"},
      {"a"},
      {"t", "s"},
      {"S"},
      {"d", "C", "P"},
   })
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Player.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

