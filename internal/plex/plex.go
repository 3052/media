package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/plex"
   "encoding/json"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
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
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/plex/user_cache.json"

   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.StringVar(&c.forwarded_for, "x", "", "x-forwarded-for")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type command struct {
   address       string
   config        maya.Config
   dash          string
   forwarded_for string
   name          string
}

func (c *command) do_address() error {
   var cache user_cache
   err := cache.User.Fetch()
   if err != nil {
      return err
   }
   address, err := plex.GetPath(c.address)
   if err != nil {
      return err
   }
   metadata, err := cache.User.RatingKey(address)
   if err != nil {
      return err
   }
   metadata, err = cache.User.Media(metadata, c.forwarded_for)
   if err != nil {
      return err
   }
   var ok bool
   cache.MediaPart, ok = metadata.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Mpd.Url, cache.Mpd.Body, err = cache.User.Mpd(
      cache.MediaPart, c.forwarded_for,
   )
   if err != nil {
      return err
   }
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd.Url, cache.Mpd.Body)
}

type user_cache struct {
   MediaPart *plex.MediaPart
   Mpd       struct {
      Body []byte
      Url  *url.URL
   }
   User plex.User
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.User.Widevine(cache.MediaPart, data)
   }
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}
