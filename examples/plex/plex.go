package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.User.Widevine(cache.MediaPart, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("plex")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.x_forwarded_for, "x", "", "x-forwarded-for")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a", "x"},
      {"d", "c", "p"},
   })
}

type command struct {
   cache maya.Cache
   // 1
   address         string
   x_forwarded_for string
   // 2
   dash string
   job  maya.WidevineJob
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
   media_part, err := metadata.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("MediaPart", media_part)
   if err != nil {
      return err
   }
   dash, err := user.Dash(media_part, c.x_forwarded_for)
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   err = c.cache.Set("User", user)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
