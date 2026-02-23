package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.Playback.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Connection *roku.Connection
   LinkCode   *roku.LinkCode
   Dash       *roku.Dash
   Playback   *roku.Playback
   User       *roku.User
}

func (c *command) do_connection() error {
   var (
      cache user_cache
      err   error
   )
   cache.Connection, err = roku.NewConnection(nil)
   if err != nil {
      return err
   }
   cache.LinkCode, err = cache.Connection.LinkCode()
   if err != nil {
      return err
   }
   fmt.Println(cache.LinkCode)
   return maya.Write(c.name, &cache)
}

func (c *command) do_set_user() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   cache.User, err = cache.Connection.User(cache.LinkCode)
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_roku() error {
   cache := &user_cache{}
   if c.get_user {
      var err error
      cache, err = maya.Read[user_cache](c.name)
      if err != nil {
         return err
      }
   }
   connection, err := roku.NewConnection(cache.User)
   if err != nil {
      return err
   }
   cache.Playback, err = connection.Playback(c.roku)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Playback.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type command struct {
   name string
   // 1
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/roku/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   // 1
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 2
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 3
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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
