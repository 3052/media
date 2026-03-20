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
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[address]:
      return c.do_address()
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {address},
      {dash_id},
   })
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
   //------------
   Job maya.Job
   //------------
   address string
   //------------
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
