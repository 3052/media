package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.Widevine,
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
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
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   if set["c"] {
      return c.do_connection()
   }
   if set["s"] {
      return c.do_set_user()
   }
   if set["r"] {
      return c.do_roku()
   }
   if set["w"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"r", "g"},
      {"w"},
      {"d"},
   })
}

func (c *client) do_connection() error {
   c.Connection = &roku.Connection{}
   err := c.Connection.Fetch(nil)
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
   if cache.Error != nil {
      return cache.Error
   }
   var err error
   c.User, err = c.Connection.User(c.LinkCode)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_roku() error {
   var user *roku.User
   if c.get_user {
      if cache.Error != nil {
         return cache.Error
      }
      user = c.User
   }
   var connection roku.Connection
   err := connection.Fetch(user)
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
   Job maya.Job
   // 5
   dash_id string
}
