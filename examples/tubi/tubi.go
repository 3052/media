package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job  maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/tubi.xml")
   if err != nil {
      return err
   }
   // 1
   flag.IntVar(&c.tubi, "t", 0, "Tubi ID")
   // 2
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
   flag.Parse()
   if c.tubi >= 1 {
      return c.do_tubi()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"t"},
      {"d", "c", "p"},
   })
}

type client struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
   // 1
   tubi int
   // 2
   dash_id string
}

///

func (c *client) do_tubi() error {
   var content tubi.Content
   err := content.Fetch(c.tubi)
   if err != nil {
      return err
   }
   var state saved_state
   state.VideoResource = &content.VideoResources[0]
   state.Dash, err = state.VideoResource.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_dash_id() error {
   var state saved_state
   err := cache.Read(&state)
   if err != nil {
      return err
   }
   job.Send = state.VideoResource.Widevine
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash_id)
}
