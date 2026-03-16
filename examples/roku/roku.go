package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Connection *roku.Connection
   Dash       *roku.Dash
   LinkCode   *roku.LinkCode
   Playback   *roku.Playback
   User       *roku.User
   // 1
   Job maya.Job
   // 2
   connection bool
   // 3
   set_user bool
   // 4
   roku     string
   get_user bool
   // 5
   dash_id string
}

///

func (c *client) do() error {
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 2
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 3
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 4
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   if set["w"] {
      return cache.Write(c)
   }
   if set["c"] {
      return c.do_connection()
   }
   if set["s"] {
      if err != nil {
         return err
      }
      return c.do_set_user()
   }
   if set["d"] {
      if err != nil {
         return err
      }
      return c.do_dash_id()
   }
   if set["r"] {
      if set["g"] {
         if err != nil {
            return err
         }
      }
      return c.do_roku()
   }
   return maya.Usage([][]string{
      {"w"},
      {"c"},
      {"s"},
      {"r", "g"},
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
func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.Widevine,
   )
}
