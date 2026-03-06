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

var cache maya.Cache

var job  maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/kanopy.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

type client struct {
   Dash         *kanopy.Dash
   Login        *kanopy.Login
   PlayManifest *kanopy.PlayManifest
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}

///

func (c *client) do_email_password() error {
   var login kanopy.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(saved_state{Login: &login})
}

func (c *client) do_address() error {
   video := &kanopy.Video{}
   err := video.Parse(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = cache.Update(&state, func() error {
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

func (c *client) do_dash_id() error {
   var state saved_state
   err := cache.Read(&state)
   if err != nil {
      return err
   }
   job.Send = func(data []byte) ([]byte, error) {
      return state.Login.Widevine(state.PlayManifest, data)
   }
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash_id)
}
