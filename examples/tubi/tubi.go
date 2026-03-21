package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/tubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringVar(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   tubi_id := maya.IntVar(&c.tubi_id, "t", "Tubi ID")
   //------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   err = maya.SetProxy(c.Proxy, "*.mp4")
   if err != nil {
      return err
   }
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[proxy]:
      return cache.Write(c)
   case set[tubi_id]:
      return c.do_tubi()
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {proxy},
      {tubi_id},
      {dash_id},
   })
}

func (c *client) do_tubi() error {
   content, err := tubi.FetchContent(c.tubi_id)
   if err != nil {
      return err
   }
   c.VideoResource = &content.VideoResources[0]
   c.Dash, err = c.VideoResource.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.VideoResource.Widevine,
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
   //-------------------------------
   Job maya.Job
   //-------------------------------
   Proxy string
   //-------------------------------
   tubi_id int
   //-------------------------------
   dash_id string
}
