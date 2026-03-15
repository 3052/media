package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "flag"
   "log"
)

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/nbc.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   if set["a"] {
      return c.do_address()
   }
   if set["w"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"a"},
      {"w"},
      {"d"},
   })
}

type client struct {
   // 1
   address string
   // 2
   Job maya.Job
   // 3
   dash_id string
}

///

func (c *client) do_address() error {
   name, err := nbc.GetName(c.address)
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(name)
   if err != nil {
      return err
   }
   stream, err := metadata.Stream()
   if err != nil {
      return err
   }
   dash, err := stream.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *client) do_dash_id() error {
   var dash nbc.Dash
   err := cache.Read(&dash)
   if err != nil {
      return err
   }
   job.Send = nbc.Widevine
   return job.DownloadDash(dash.Body, dash.Url, c.dash_id)
}

