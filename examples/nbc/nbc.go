package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "flag"
   "log"
)

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
   err = cache.Update(c, func() error {
      c.Dash, err = stream.Dash()
      return err
   }, true)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
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
   flag.IntVar(&c.Job.Threads, "t", 2, "threads")
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
      {"d", "t"},
   })
}

type client struct {
   Dash *nbc.Dash
   // 1
   address string
   // 2
   Job maya.Job
   // 3
   dash_id string
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, nbc.Widevine)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
