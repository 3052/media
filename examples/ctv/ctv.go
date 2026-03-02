package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4a,*.m4v")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/ctv.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
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
      {"a"},
      {"d", "c", "p"},
   })
}

func (c *client) do_address() error {
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
   err = c.cache.Set(dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
func (c *client) do_dash() error {
   c.job.Send = ctv.Widevine
   var dash ctv.Dash
   err := c.cache.Get(&dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

type client struct {
   cache maya.Cache
   // 1
   address string
   // 2
   dash string
   job  maya.WidevineJob
}

