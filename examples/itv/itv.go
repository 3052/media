package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

type saved_state struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
}

func (c *client) do_playlist() error {
   var title itv.Title
   title.LatestAvailableVersion.PlaylistUrl = c.playlist
   playlist, err := title.Playlist()
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Update(&state, func() error {
      state.MediaFile, err = playlist.FullHd()
      if err != nil {
         return err
      }
      state.Dash, err = state.MediaFile.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func main() {
   // ALL REQUEST ARE GEO BLOCKED
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_address() error {
   titles, err := itv.Titles(itv.LegacyId(c.address))
   if err != nil {
      return err
   }
   for i, title := range titles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&title)
   }
   return nil
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.MediaFile.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type client struct {
   cache maya.Cache
   // 1
   address string
   // 2
   playlist string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/itv.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.playlist != "" {
      return c.do_playlist()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a"},
      {"p"},
      {"d", "C", "P"},
   })
}
