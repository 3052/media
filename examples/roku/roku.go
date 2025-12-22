package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/roku"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.Widevine(data)
   }
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

type user_cache struct {
   Connection *roku.Connection
   LinkCode   *roku.LinkCode
   Mpd        *roku.Mpd
   Playback   *roku.Playback
   User       *roku.User
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/roku/userCache.xml"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.BoolVar(&c.connection, "c", false, "connection")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.set_user, "s", false, "set user")
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
   flag.Usage()
   return nil
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   return write(c.name, &cache)
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func (c *command) do_set_user() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.User, err = cache.Connection.User(cache.LinkCode)
   if err != nil {
      return err
   }
   return write(c.name, cache)
}

func (c *command) do_roku() error {
   cache := &user_cache{}
   if c.get_user {
      var err error
      cache, err = read(c.name)
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
   cache.Mpd, err = cache.Playback.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd.Url, cache.Mpd.Body)
}

type command struct {
   name   string
   config maya.Config
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
