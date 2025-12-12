package main

import (
   "41.neocities.org/media/itv"
   "41.neocities.org/net"
   "encoding/json"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func main() {
   // ALL REQUEST ARE GEO BLOCKED
   net.Transport(func(req *http.Request) string {
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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/itv/user_cache.json"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   flag.IntVar(&c.config.Threads, "t", 2, "threads")
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
   flag.Usage()
   return nil
}

func (c *command) do_address() error {
   legacy_id, err := itv.LegacyId(c.address)
   if err != nil {
      return err
   }
   titles, err := itv.Titles(legacy_id)
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

type command struct {
   address  string
   config   net.Config
   dash     string
   name     string
   playlist string
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
   cache.Mpd, cache.MpdBody, err = cache.MediaFile.Mpd()
   if err != nil {
      return err
   }
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd, cache.MpdBody)
}

type user_cache struct {
   MediaFile *itv.MediaFile
   Mpd       *url.URL
   MpdBody   []byte
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.MediaFile.Widevine(data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}
