package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "flag"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/kanopy.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

type saved_state struct {
   Dash         *kanopy.Dash
   Login        *kanopy.Login
   PlayManifest *kanopy.PlayManifest
}

func (c *client) do_email_password() error {
   var login kanopy.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Login: &login})
}

type client struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *client) do_address() error {
   video := &kanopy.Video{}
   err := video.Parse(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Update(&state, func() error {
      if video.VideoId == 0 {
         video, err = state.Login.Video(video.Alias)
         if err != nil {
            return err
         }
      }
      member, err := state.Login.Membership()
      if err != nil {
         return err
      }
      plays, err := state.Login.Plays(member, video.VideoId)
      if err != nil {
         return err
      }
      for _, caption := range plays.Captions {
         for _, file := range caption.Files {
            fmt.Println(file.Url)
         }
      }
      state.PlayManifest, err = plays.Dash()
      if err != nil {
         return err
      }
      state.Dash, err = state.PlayManifest.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return state.Login.Widevine(state.PlayManifest, data)
   }
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

