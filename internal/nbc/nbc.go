package main

import (
   "41.neocities.org/media/nbc"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
      return "LP"
   })
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
   c.name = cache + "/nbc/mpd.json"

   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.nbc, "n", 0, "NBC ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.Parse()
   if c.nbc >= 1 {
      return c.do_nbc()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type command struct {
   name   string
   config  net.Config
   dash string
   nbc     int
}

type mpd struct {
   Body []byte
   Url *url.URL
}

func (c *command) do_nbc() error {
   metadata, err := nbc.FetchMetadata(c.nbc)
   if err != nil {
      return err
   }
   stream_info, err := metadata.StreamInfo()
   if err != nil {
      return err
   }
   var cache mpd
   cache.Url, cache.Body, err = stream_info.Mpd()
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
   return net.Representations(cache.Url, cache.Body)
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache mpd
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = nbc.Widevine
   return c.config.Download(cache.Url, cache.Body, c.dash)
}
