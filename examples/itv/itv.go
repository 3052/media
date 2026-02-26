package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rosso/itv.xml"
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
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

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.MediaFile.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
}

func main() {
   // ALL REQUEST ARE GEO BLOCKED
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_address() error {
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

func (c *command) do_playlist() error {
   var title itv.Title
   title.LatestAvailableVersion.PlaylistUrl = c.playlist
   playlist, err := title.Playlist()
   if err != nil {
      return err
   }
   var cache user_cache
   cache.MediaFile, err = playlist.FullHd()
   if err != nil {
      return err
   }
   cache.Dash, err = cache.MediaFile.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type command struct {
   name string
   // 1
   address string
   // 2
   playlist string
   // 3
   dash string
   job  maya.WidevineJob
}
