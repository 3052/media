package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/itv"
   "errors"
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
   c.name = cache + "/itv/userCache.xml"
   c.job.ClientId = filepath.Join(cache, "/L3/client_id.bin")
   c.job.PrivateKey = filepath.Join(cache, "/L3/private_key.pem")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   // 1
   if c.address != "" {
      return c.do_address()
   }
   // 2
   if c.playlist != "" {
      return c.do_playlist()
   }
   // 3
   if c.dash != "" {
      return c.do_dash()
   }
   maya.Usage([][]string{
      {"a"},
      {"p"},
      {"d", "C", "P"},
   })
   return nil
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.MediaFile.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
}

func main() {
   // ALL REQUEST ARE GEO BLOCKED
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".dash" {
         return ""
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
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
   var (
      cache user_cache
      ok    bool
   )
   cache.MediaFile, ok = playlist.FullHd()
   if !ok {
      return errors.New(".FullHd()")
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
