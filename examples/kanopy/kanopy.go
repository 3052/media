package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/kanopy"
   "encoding/xml"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Login.Widevine(cache.StreamInfo, data)
   }
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

type user_cache struct {
   Login      *kanopy.Login
   Mpd        *kanopy.Mpd
   StreamInfo *kanopy.StreamInfo
}

func main() {
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
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
   c.name = cache + "/kanopy/userCache.xml"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.IntVar(&c.kanopy, "k", 0, "Kanopy ID")
   flag.StringVar(&c.password, "p", "", "password")
   flag.Parse()

   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.kanopy >= 1 {
      return c.do_kanopy()
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

func (c *command) do_email_password() error {
   var cache user_cache
   err := cache.Login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
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

type command struct {
   config   maya.Config
   dash     string
   email    string
   kanopy   int
   name     string
   password string
}

func (c *command) do_kanopy() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   member, err := cache.Login.Membership()
   if err != nil {
      return err
   }
   play, err := cache.Login.PlayResponse(member, c.kanopy)
   if err != nil {
      return err
   }
   var ok bool
   cache.StreamInfo, ok = play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Mpd, err = cache.StreamInfo.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd.Url, cache.Mpd.Body)
}
