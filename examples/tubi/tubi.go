package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rosso/tubi.xml"
   // 1
   flag.IntVar(&c.tubi, "t", 0, "Tubi ID")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "T", 2, "threads")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".mp4" {
         return "", false
      }
      return c.proxy, true
   })
   if c.tubi >= 1 {
      return c.do_tubi()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"t", "x"},
      {"d", "T", "c", "p"},
   })
}

func (c *command) do_tubi() error {
   var content tubi.Content
   err := content.Fetch(c.tubi)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.VideoResource = &content.VideoResources[0]
   cache.Dash, err = cache.VideoResource.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.VideoResource.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type user_cache struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
}

type command struct {
   name string
   // 1
   tubi  int
   proxy string
   // 2
   dash string
   job  maya.WidevineJob
}
