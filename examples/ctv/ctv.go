package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("ctv")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".m4a", ".m4v":
         return "", false
      }
      return c.proxy, true
   })
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

func (c *command) do_address() error {
   link_path, err := ctv.GetPath(c.address)
   if err != nil {
      return err
   }
   resolve, err := ctv.Resolve(link_path)
   if err != nil {
      return err
   }
   axis, err := resolve.AxisContent()
   if err != nil {
      return err
   }
   playback, err := axis.Playback()
   if err != nil {
      return err
   }
   manifest, err := axis.Manifest(playback)
   if err != nil {
      return err
   }
   dash, err := manifest.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_dash() error {
   c.job.Send = ctv.Widevine
   var dash ctv.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

type command struct {
   cache maya.Cache
   // 1
   address string
   proxy   string
   // 2
   dash string
   job  maya.WidevineJob
}
