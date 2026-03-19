package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "flag"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/nbc.xml")
   if err != nil {
      return err
   }
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   address := maya.StringVar(&c.address, "a", "address")
   // 3
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {address},
         {dash_id},
      })
   case set[widevine]:
      return cache.Write(c)
   case read_err != nil:
      return read_err
   case set[address]:
      return c.do_address()
   case set[dash_id]:
      return c.do_dash_id()
   }
   return nil
}

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
   c.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Dash *nbc.Dash
   // 1
   Job maya.Job
   // 2
   address string
   // 3
   dash_id string
}

func (c *client) do_dash_id() error {
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
