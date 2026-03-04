package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   // 1
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 2
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 3
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
   flag.Parse()
   if c.connection {
      return c.do_connection()
   }
   if c.set_user {
      return c.do_set_user()
   }
   if c.roku != "" {
      return c.do_roku()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"r", "g"},
      {"d", "C", "P"},
   })
}

var cache maya.Cache

var job  maya.WidevineJob

type client struct {
   // 1
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   // 4
   dash string
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var state struct {
   Dash       *roku.Dash
   Connection *roku.Connection
   LinkCode   *roku.LinkCode
   Playback   *roku.Playback
   User       *roku.User
}

func (c *client) do_connection() error {
   var err error
   state.Connection, err = roku.NewConnection(nil)
   if err != nil {
      return err
   }
   state.LinkCode, err = state.Connection.LinkCode()
   if err != nil {
      return err
   }
   fmt.Println(state.LinkCode)
   return cache.Write(state)
}

func (c *client) do_roku() error {
   var user *roku.User
   if c.get_user {
      _, err := cache.Read(&state)
      if err != nil {
         return err
      }
      user = state.User
   }
   connection, err := roku.NewConnection(user)
   if err != nil {
      return err
   }
   state.Playback, err = connection.Playback(c.roku)
   if err != nil {
      return err
   }
   state.Dash, err = state.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_dash() error {
   _, err := cache.Read(&state)
   if err != nil {
      return err
   }
   job.Send = state.Playback.Widevine
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do_set_user() error {
   _, err := cache.Read(&state)
   if err != nil {
      return err
   }
   state.User, err = state.Connection.User(state.LinkCode)
   if err != nil {
      return err
   }
   return cache.Write(state)
}
