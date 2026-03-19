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
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   tubi_id := maya.IntVar(&c.tubi_id, "t", "Tubi ID")
   // 3
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {tubi_id},
         {dash_id},
      })
   case set[widevine]:
      return cache.Write(c)
   case set[tubi_id]:
      return c.do_tubi()
   case read_err != nil:
      return read_err
   case set[dash_id]:
      return c.do_dash_id()
   }
   return nil
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.VideoResource.Widevine,
   )
}

type client struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
   // 1
   Job maya.Job
   // 2
   tubi_id int
   // 3
   dash_id string
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
