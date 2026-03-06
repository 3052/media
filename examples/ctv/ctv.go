package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "flag"
   "log"
)

func (c *client) do_dash() error {
   var dash ctv.Dash
   err := cache.Read(&dash)
   if err != nil {
      return err
   }
   job.Send = ctv.Widevine
   return job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4a,*.m4v")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/ctv.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
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

type client struct {
   // 1
   address string
   // 2
   dash string
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
   err = cache.Write(dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
