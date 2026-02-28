package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) do_set_user() error {
   var connection roku.Connection
   err := c.cache.Get("Connection", &connection)
   if err != nil {
      return err
   }
   var link_code roku.LinkCode
   err = c.cache.Get("LinkCode", &link_code)
   if err != nil {
      return err
   }
   user, err := connection.User(&link_code)
   if err != nil {
      return err
   }
   return c.cache.Set("User", user)
}

func (c *command) do_roku() error {
   var user *roku.User
   if c.get_user {
      user = &roku.User{}
      err := c.cache.Get("User", user)
      if err != nil {
         return err
      }
   }
   connection, err := roku.NewConnection(user)
   if err != nil {
      return err
   }
   playback, err := connection.Playback(c.roku)
   if err != nil {
      return err
   }
   err = c.cache.Set("Playback", playback)
   if err != nil {
      return err
   }
   dash, err := playback.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
func (c *command) do_dash() error {
   var dash roku.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   var playback roku.Playback
   err = c.cache.Get("Playback", &playback)
   if err != nil {
      return err
   }
   c.job.Send = playback.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
   // 1
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   proxy    string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("roku")
   // 1
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 2
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 3
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return c.proxy, path.Ext(req.URL.Path) != ".mp4"
   })
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
func (c *command) do_connection() error {
   connection, err := roku.NewConnection(nil)
   if err != nil {
      return err
   }
   err = c.cache.Set("Connection", connection)
   if err != nil {
      return err
   }
   link_code, err := connection.LinkCode()
   if err != nil {
      return err
   }
   fmt.Println(link_code)
   return c.cache.Set("LinkCode", link_code)
}
