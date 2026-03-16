package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
)

func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.VideoResource.Widevine,
   )
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

func (c *client) do() error {
   err := cache.Setup("rosso/tubi.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
   // 1
   flag.IntVar(&c.tubi, "t", 0, "Tubi ID")
   // 2
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   if set["t"] {
      return c.do_tubi()
   }
   if set["w"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"t"},
      {"w"},
      {"d"},
   })
}

type client struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
   // 1
   tubi int
   // 2
   Job maya.Job
   // 3
   dash_id string
}

func (c *client) do_tubi() error {
   var content tubi.Content
   err := content.Fetch(c.tubi)
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
