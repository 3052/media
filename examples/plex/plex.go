package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/plex"
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
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.User.Widevine(cache.MediaPart, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
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

type user_cache struct {
   Dash       *plex.Dash
   MediaPart *plex.MediaPart
   User      *plex.User
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   c.name = cache + "/plex/userCache.xml"
   c.job.ClientId = filepath.Join(cache, "/L3/client_id.bin")
   c.job.PrivateKey = filepath.Join(cache, "/L3/private_key.pem")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.x_forwarded_for, "x", "", "x-forwarded-for")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   // 1
   if c.address != "" {
      return c.do_address()
   }
   // 2
   if c.dash != "" {
      return c.do_dash()
   }
   maya.Usage([][]string{
      {"a", "x"},
      {"d", "c", "p"},
   })
   return nil
}

type command struct {
   name            string
   // 1
   address         string
   x_forwarded_for string
   // 2
   dash            string
   job          maya.WidevineJob
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
   cache.Dash, err = user.Dash(cache.MediaPart, c.x_forwarded_for)
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
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}
