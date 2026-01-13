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
   "path"
   "path/filepath"
)

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

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
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
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/roku/userCache.xml"
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   // 1
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 2
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 3
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.Parse()
   // 1
   if c.connection {
      return c.do_connection()
   }
   // 2
   if c.set_user {
      return c.do_set_user()
   }
   // 3
   if c.roku != "" {
      return c.do_roku()
   }
   // 4
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type command struct {
   job     maya.WidevineJob
   name       string
   
   // 1
   connection bool
   // 2
   set_user   bool
   // 3
   roku       string
   get_user   bool
   // 4
   dash       string
}

///

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

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.Widevine(data)
   }
   return c.job.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}
