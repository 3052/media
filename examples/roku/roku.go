package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Playback.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

var cache maya.Cache

var job maya.WidevineJob

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

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
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
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
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"r", "g"},
      {"d", "C", "P"},
   })
}

func (c *client) do_connection() error {
   var err error
   c.Connection, err = roku.NewConnection(nil)
   if err != nil {
      return err
   }
   c.LinkCode, err = c.Connection.LinkCode()
   if err != nil {
      return err
   }
   fmt.Println(c.LinkCode)
   return cache.Write(c)
}

func (c *client) do_set_user() error {
   return cache.Update(c, func() error {
      var err error
      c.User, err = c.Connection.User(c.LinkCode)
      return err
   })
}

type client struct {
   Connection *roku.Connection
   Dash       *roku.Dash
   LinkCode   *roku.LinkCode
   Playback   *roku.Playback
   User       *roku.User
   // 1
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   // 4
   dash_id string
}

func (c *client) do_roku() error {
   var user *roku.User
   if c.get_user {
      err := cache.Read(c)
      if err != nil {
         return err
      }
      user = c.User
   }
   connection, err := roku.NewConnection(user)
   if err != nil {
      return err
   }
   c.Playback, err = connection.Playback(c.roku)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
