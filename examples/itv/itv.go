package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.MediaFile.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

func main() {
   log.SetFlags(log.Ltime)
   // ALL REQUEST ARE GEO BLOCKED
   maya.SetProxy("", "*.dash")
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
   err := cache.Setup("rosso/itv.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.playlist != "" {
      return c.do_playlist()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"a"},
      {"p"},
      {"d", "C", "P"},
   })
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

type client struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
   // 1
   address string
   // 2
   playlist string
   // 3
   dash_id string
}

func (c *client) do_playlist() error {
   var title itv.Title
   title.LatestAvailableVersion.PlaylistUrl = c.playlist
   playlist, err := title.Playlist()
   if err != nil {
      return err
   }
   err = cache.Update(c, func() error {
      c.MediaFile, err = playlist.FullHd()
      if err != nil {
         return err
      }
      c.Dash, err = c.MediaFile.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
