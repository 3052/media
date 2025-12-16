package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/plex"
   "encoding/xml"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func do_ifconfig() error {
   resp, err := http.Get("http://ifconfig.co")
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = os.Stdout.ReadFrom(resp.Body)
   if err != nil {
      return err
   }
   return nil
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.User.Widevine(cache.MediaPart, data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}

type user_cache struct {
   MediaPart *plex.MediaPart
   Mpd       *url.URL
   MpdBody   []byte
   User      *plex.User
}

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
   c.name = cache + "/plex/userCache.xml"

   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")

   flag.BoolVar(&c.ifconfig, "i", false, "ifconfig.co")

   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.x_forwarded_for, "x", "", "x-forwarded-for")

   flag.StringVar(&c.dash, "d", "", "DASH ID")

   flag.Parse()

   if c.ifconfig {
      return do_ifconfig()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

func (c *command) do_address() error {
   var user plex.User
   err := user.Fetch()
   if err != nil {
      return err
   }
   address, err := plex.GetPath(c.address)
   if err != nil {
      return err
   }
   metadata, err := user.RatingKey(address)
   if err != nil {
      return err
   }
   metadata, err = user.Media(metadata, c.x_forwarded_for)
   if err != nil {
      return err
   }
   var (
      cache user_cache
      ok    bool
   )
   cache.MediaPart, ok = metadata.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Mpd, cache.MpdBody, err = user.Mpd(cache.MediaPart, c.x_forwarded_for)
   if err != nil {
      return err
   }
   cache.User = &user
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd, cache.MpdBody)
}

type command struct {
   address         string
   config          maya.Config
   dash            string
   ifconfig        bool
   name            string
   x_forwarded_for string
}
