package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) do_dash() error {
   var dash tubi.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   var video tubi.VideoResource
   err = c.cache.Get("VideoResource", &video)
   if err != nil {
      return err
   }
   c.job.Send = video.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
   // 1
   tubi  int
   // 2
   dash string
   job  maya.WidevineJob
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("tubi")
   // 1
   flag.IntVar(&c.tubi, "t", 0, "Tubi ID")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.tubi >= 1 {
      return c.do_tubi()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"t", "x"},
      {"d", "c", "p"},
   })
}

func (c *command) do_tubi() error {
   var content tubi.Content
   err := content.Fetch(c.tubi)
   if err != nil {
      return err
   }
   video := content.VideoResources[0]
   err = c.cache.Set("VideoResource", video)
   if err != nil {
      return err
   }
   dash, err := video.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
